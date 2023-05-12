// repository.user.go

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	pb "commonpb"
)

func (s *server) IsUserValid(ctx context.Context, args *pb.AccountRequest) (*pb.BoolReply, error) {
	fmt.Println("IsUserValid")

	raftNode := getRaftNode()

	username := args.GetUsername()
	password := args.GetPassword()

	raftNode.userListRWMutex.RLock()
	defer raftNode.userListRWMutex.RUnlock()

	for _, u := range raftNode.userList {
		if u.Username == username && u.Password == password {
			return &pb.BoolReply{Reply: true}, nil
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &pb.BoolReply{Reply: false}, nil
	}
}

func (s *server) RegisterNewUser(ctx context.Context, args *pb.AccountRequest) (*pb.BoolReply, error) {
	fmt.Println("RegisterNewUser")

	username := args.GetUsername()
	password := args.GetPassword()

	if strings.TrimSpace(password) == "" {
		return &pb.BoolReply{Reply: false}, errors.New("The password can't be empty")
	} else if !isUsernameAvailable2(username) {
		return &pb.BoolReply{Reply: false}, errors.New("The username isn't available")
	}

	up := userProposal{
		Username: username,
		Password: password,
	}
	bytes, _ := json.Marshal(up)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		proposal := []byte("addUser:" + string(bytes))
		err := commitProposal(proposal)
		if err != nil {
			return nil, err
		}
		return &pb.BoolReply{Reply: true}, nil
	}
}

func (s *server) IsUsernameAvailable(ctx context.Context, args *pb.UsernameRequest) (*pb.BoolReply, error) {
	fmt.Println("IsUsernameAvailable")

	raftNode := getRaftNode()

	username := args.GetUsername()

	raftNode.userListRWMutex.RLock()
	defer raftNode.userListRWMutex.RUnlock()

	for _, u := range raftNode.userList {
		if u.Username == username {
			return &pb.BoolReply{Reply: false}, nil
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &pb.BoolReply{Reply: true}, nil
	}
}

func (s *server) IsFollowed(ctx context.Context, args *pb.FollowRequest) (*pb.IsFollowedReply, error) {
	fmt.Println("IsFollowed")

	raftNode := getRaftNode()

	thisUser := args.GetThisUser()
	followUser := args.GetFollowUser()

	if isUsernameAvailable2(followUser) {
		return &pb.IsFollowedReply{Found: false, Index: -1}, errors.New("The user doesn't exist")
	}

	usernames := raftNode.follows[thisUser]
	usernames.RLock()
	defer usernames.RUnlock()

	for i, a := range usernames.followList {
		if a == followUser {
			return &pb.IsFollowedReply{Found: true, Index: int64(i)}, nil
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &pb.IsFollowedReply{Found: false, Index: -1}, nil
	}
}

func (s *server) AddFollow(ctx context.Context, args *pb.FollowRequest) (*pb.Reply, error) {
	fmt.Println("AddFollow")

	followUser := args.GetFollowUser()
	thisUser := args.GetThisUser()

	found, _, err := isFollowed2(thisUser, followUser)

	if err != nil {
		return &pb.Reply{Message: "error"}, err
	}

	if found {
		return &pb.Reply{Message: "error"}, errors.New("This user is already followed")
	}

	fp := followAddProposal{
		ThisUser:   thisUser,
		FollowUser: followUser,
	}
	bytes, _ := json.Marshal(fp)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		proposal := []byte("addFollow:" + string(bytes))
		err = commitProposal(proposal)
		if err != nil {
			return nil, err
		}
		return &pb.Reply{Message: "success"}, nil
	}
}

func (s *server) RemoveFollow(ctx context.Context, args *pb.FollowRequest) (*pb.Reply, error) {
	fmt.Println("RemoveFollow")

	thisUser := args.GetThisUser()
	followUser := args.GetFollowUser()

	found, i, err := isFollowed2(thisUser, followUser)

	if err != nil {
		return &pb.Reply{Message: "error"}, err
	}

	if !found {
		return &pb.Reply{Message: "error"}, errors.New("This user isn't followed")
	}

	fp := followRemoveProposal{
		ThisUser:    thisUser,
		RemoveIndex: i,
	}
	bytes, _ := json.Marshal(fp)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		proposal := []byte("removeFollow:" + string(bytes))
		err = commitProposal(proposal)
		if err != nil {
			return nil, err
		}
		return &pb.Reply{Message: "success"}, nil
	}
}

func (s *server) GetFollowedUsers(ctx context.Context, args *pb.UsernameRequest) (*pb.UsernameListReply, error) {
	fmt.Println("GetFollowedUsers")

	raftNode := getRaftNode()

	thisUser := args.GetUsername()

	if thisUser == "" {
		return &pb.UsernameListReply{FollowList: []string{}}, errors.New("The user of this account not found")
	}

	f := raftNode.follows[thisUser]
	f.RLock()
	defer f.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &pb.UsernameListReply{FollowList: f.followList}, nil
	}
}

/*
	Non-RPC versions for use within backend functions
*/
func isUsernameAvailable2(username string) bool {
	fmt.Println("isUsernameAvailable2")

	raftNode := getRaftNode()

	raftNode.userListRWMutex.RLock()
	defer raftNode.userListRWMutex.RUnlock()

	for _, u := range raftNode.userList {
		if u.Username == username {
			return false
		}
	}
	return true
}

func isFollowed2(thisUser string, followUser string) (bool, int, error) {
	fmt.Println("isFollowed2")

	raftNode := getRaftNode()

	if isUsernameAvailable2(followUser) {
		return false, -1, errors.New("The user doesn't exist")
	}

	usernames := raftNode.follows[thisUser]
	usernames.RLock()
	defer usernames.RUnlock()

	for i, a := range usernames.followList {
		if a == followUser {
			return true, i, nil
		}
	}

	return false, -1, nil
}

func getFollowedUsers2(username string) ([]string, error) {
	fmt.Println("getFollowedUsers2")

	raftNode := getRaftNode()

	if username == "" {
		return nil, errors.New("The user of this account not found")
	}

	f := raftNode.follows[username]
	f.RLock()
	defer f.RUnlock()

	return f.followList, nil
}

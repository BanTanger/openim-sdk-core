package group

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/syncdb"
	"open_im_sdk/pkg/utils"
)

func (g *Group) SyncGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	resp, err := util.CallApi[group.GetGroupMembersInfoResp](ctx, constant.GetGroupMembersInfoRouter, &group.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs})
	if err != nil {
		return err
	}
	var members []any
	for _, member := range resp.Members {
		members = append(members, &model_struct.LocalGroupMember{
			GroupID:        member.GroupID,
			UserID:         member.UserID,
			Nickname:       member.Nickname,
			FaceURL:        member.FaceURL,
			RoleLevel:      member.RoleLevel,
			JoinTime:       member.JoinTime,
			JoinSource:     member.JoinSource,
			InviterUserID:  member.InviterUserID,
			MuteEndTime:    member.MuteEndTime,
			OperatorUserID: member.OperatorUserID,
			Ex:             member.Ex,
			//AttachedInfo:   member.AttachedInfo, // todo
		})
	}
	return syncdb.NewSync(g.db.GetDB(ctx)).AddChange(members).Start()
}

func (g *Group) SyncGroup(ctx context.Context, groupID string) error {
	resp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: []string{groupID}})
	if err != nil {
		return err
	}
	if len(resp.GroupInfos) == 0 {
		return errs.ErrGroupIDNotFound.Wrap(groupID)
	}
	groupInfo := resp.GroupInfos[0]
	groupModel := &model_struct.LocalGroup{
		GroupID:                groupInfo.GroupID,
		GroupName:              groupInfo.GroupName,
		Notification:           groupInfo.Notification,
		Introduction:           groupInfo.Introduction,
		FaceURL:                groupInfo.FaceURL,
		CreateTime:             groupInfo.CreateTime,
		Status:                 groupInfo.Status,
		CreatorUserID:          groupInfo.CreatorUserID,
		GroupType:              groupInfo.GroupType,
		OwnerUserID:            groupInfo.OwnerUserID,
		MemberCount:            int32(groupInfo.MemberCount),
		Ex:                     groupInfo.Ex,
		NeedVerification:       groupInfo.NeedVerification,
		LookMemberInfo:         groupInfo.LookMemberInfo,
		ApplyMemberFriend:      groupInfo.ApplyMemberFriend,
		NotificationUpdateTime: groupInfo.NotificationUpdateTime,
		NotificationUserID:     groupInfo.NotificationUserID,
		//AttachedInfo:           groupInfo.AttachedInfo, // TODO
	}
	if err := syncdb.NewSync(g.db.GetDB(ctx)).AddChange(groupModel).Start(); err != nil {
		return err
	}
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupModel))
	return nil
}

func (g *Group) SyncGroupAndMember(ctx context.Context, groupID string) error {
	groupResp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: []string{groupID}})
	if err != nil {
		return err
	}
	if len(groupResp.GroupInfos) == 0 {
		return errs.ErrGroupIDNotFound.Wrap(groupID)
	}
	groupInfo := groupResp.GroupInfos[0]
	showNumber := int32(20)
	var members []any
	for i := int32(0); ; i++ {
		memberReq := &group.GetGroupMemberListReq{GroupID: groupInfo.GroupID, Pagination: &sdkws.RequestPagination{PageNumber: i, ShowNumber: showNumber}}
		memberResp, err := util.CallApi[group.GetGroupMemberListResp](ctx, constant.GetGroupAllMemberListRouter, memberReq)
		if err != nil {
			return err
		}
		for _, member := range memberResp.Members {
			members = append(members, &model_struct.LocalGroupMember{
				GroupID:        member.GroupID,
				UserID:         member.UserID,
				Nickname:       member.Nickname,
				FaceURL:        member.FaceURL,
				RoleLevel:      member.RoleLevel,
				JoinTime:       member.JoinTime,
				JoinSource:     member.JoinSource,
				InviterUserID:  member.InviterUserID,
				MuteEndTime:    member.MuteEndTime,
				OperatorUserID: member.OperatorUserID,
				Ex:             member.Ex,
				//AttachedInfo:   member.AttachedInfo, // todo
			})
		}
		if int32(len(memberResp.Members)) < showNumber {
			break
		}
	}
	groupModel := &model_struct.LocalGroup{
		GroupID:       groupInfo.GroupID,
		GroupName:     groupInfo.GroupName,
		Notification:  groupInfo.Notification,
		Introduction:  groupInfo.Introduction,
		FaceURL:       groupInfo.FaceURL,
		CreateTime:    groupInfo.CreateTime,
		Status:        groupInfo.Status,
		CreatorUserID: groupInfo.CreatorUserID,
		GroupType:     groupInfo.GroupType,
		OwnerUserID:   groupInfo.OwnerUserID,
		MemberCount:   int32(groupInfo.MemberCount),
		Ex:            groupInfo.Ex,
		//AttachedInfo:           groupInfo.AttachedInfo, // TODO
		NeedVerification:       groupInfo.NeedVerification,
		LookMemberInfo:         groupInfo.LookMemberInfo,
		ApplyMemberFriend:      groupInfo.ApplyMemberFriend,
		NotificationUpdateTime: groupInfo.NotificationUpdateTime,
		NotificationUserID:     groupInfo.NotificationUserID,
	}
	s := syncdb.NewSync(g.db.GetDB(ctx)).AddChange(groupModel).AddComplete([]string{"group_id"}, members)
	if err := s.Start(); err != nil {
		return err
	}
	g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupModel))
	for _, member := range members {
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(member))
	}
	return nil
}

func (g *Group) SyncSelfGroupApplication(ctx context.Context) error {
	list, err := GetAll(ctx, constant.GetSendGroupApplicationListRouter, &group.GetUserReqApplicationListReq{}, func(resp *group.GetGroupApplicationListResp) []*sdkws.GroupRequest { return resp.GroupRequests })
	if err != nil {
		return err
	}

	svrList, err := g.getSendGroupApplicationListFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getSendGroupApplicationListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalSendGroupRequest(svrList)
	onLocal, err := g.db.GetSendGroupApplication()
	if err != nil {
		log.NewError(operationID, "GetSendGroupApplication failed ", err.Error())
		return
	}

	log.NewInfo(operationID, "svrList onServer onLocal ", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupRequestDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupRequest failed ", err.Error(), *onServer[index])
			continue
		}
		callbackData := *onServer[index]
		if g.listener != nil {
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAdded ", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == constant.GroupResponseRefuse {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))
			}

		} else if onServer[index].HandleResult == constant.GroupResponseAgree {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
			if g.listenerForService != nil {
				g.listenerForService.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
		} else {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := *onLocal[index]
		if g.listener != nil {
			g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
		}
		log.Info(operationID, "OnGroupApplicationDeleted", utils.StructToJsonString(callbackData))
	}
}

func GetAll[A interface {
	GetPagination() *sdkws.RequestPagination
}, B, C any](ctx context.Context, router string, req A, fn func(resp *B) []C) ([]C, error) {
	if req.GetPagination().ShowNumber == 0 {
		req.GetPagination().ShowNumber = 50
	}
	var res []C
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i
		memberResp, err := util.CallApi[B](ctx, router, req)
		if err != nil {
			return nil, err
		}
		list := fn(memberResp)
		res = append(res, list...)
		if len(list) < int(req.GetPagination().ShowNumber) {
			break
		}
	}
	return res, nil
}
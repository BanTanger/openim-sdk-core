package super_group

import (
	"context"
	"open_im_sdk/internal/util"
)

func (s *SuperGroup) SyncJoinedGroupList(ctx context.Context) error {
	list, err := s.getJoinedGroupListFromSvr(ctx)
	if err != nil {
		return err
	}
	localData, err := s.db.GetJoinedSuperGroupList(ctx)
	if err != nil {
		return err
	}
	return s.syncerGroup.Sync(ctx, util.Batch(ServerGroupToLocalGroup, list), localData, nil)
}

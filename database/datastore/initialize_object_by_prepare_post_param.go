package datastore_db

import (
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func InitializeObjectByPreparePostParam(ownerPID types.PID, param datastore_types.DataStorePreparePostParamV1) (uint32, *nex.Error) {
	var dataID uint32

	tagArray := make([]string, 0, len(param.Tags))
	for i := range param.Tags {
		tagArray = append(tagArray, string(param.Tags[i]))
	}

	now := time.Now()
	err := database.Postgres.QueryRow(`INSERT INTO datastore.objects (
		owner,
		size,
		name,
		data_type,
		meta_binary,
		permission,
		permission_recipients,
		delete_permission,
		delete_permission_recipients,
		flag,
		period,
		refer_data_id,
		tags,
		creation_date,
		update_date
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15
	) RETURNING data_id`,
		ownerPID,
		param.Size,
		param.Name,
		param.DataType,
		param.MetaBinary,
		param.Permission.Permission,
		pq.Array(param.Permission.RecipientIDs),
		param.DelPermission.Permission,
		pq.Array(param.DelPermission.RecipientIDs),
		param.Flag,
		param.Period,
		param.ReferDataID,
		pq.Array(tagArray),
		now,
		now,
	).Scan(&dataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return 0, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return dataID, nil
}

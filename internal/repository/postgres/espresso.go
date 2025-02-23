// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/espresso/table"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-jet/jet/v2/postgres"
)

func (r *postgresRepository) GetEspressoNonce(
	ctx context.Context,
	senderAddress string,
	nameOrAddress string,
) (uint64, error) {
	// assume all are hex address string
	sender := common.HexToAddress(senderAddress)
	app := common.HexToAddress(nameOrAddress)
	sel := table.EspressoNonce.
		SELECT(table.EspressoNonce.Nonce).
		FROM(table.EspressoNonce).
		WHERE(
			table.EspressoNonce.SenderAddress.EQ(postgres.Bytea(sender.Bytes())).
				AND(table.EspressoNonce.ApplicationAddress.EQ(postgres.Bytea(app.Bytes()))),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var nonce uint64
	err := row.Scan(
		&nonce,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (r *postgresRepository) UpdateEspressoNonce(
	ctx context.Context,
	senderAddress string,
	nameOrAddress string,
) error {
	// assume all are hex address string
	sender := common.HexToAddress(senderAddress)
	app := common.HexToAddress(nameOrAddress)
	nonce, err := r.GetEspressoNonce(ctx, senderAddress, nameOrAddress)
	if err != nil {
		return err
	}
	nextNonce := nonce + 1

	nonceInsertStmt := table.EspressoNonce.INSERT(
		table.EspressoNonce.SenderAddress,
		table.EspressoNonce.ApplicationAddress,
		table.EspressoNonce.Nonce,
	).VALUES(
		postgres.Bytea(sender.Bytes()),
		postgres.Bytea(app.Bytes()),
		nextNonce,
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	sqlStr, args := nonceInsertStmt.
		ON_CONFLICT(table.EspressoNonce.SenderAddress, table.EspressoNonce.ApplicationAddress).
		DO_UPDATE(postgres.SET(
			table.EspressoNonce.Nonce.SET(postgres.RawInt(fmt.Sprintf("%d", nextNonce))),
		)).Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)

	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return nil
}

func (r *postgresRepository) GetInputIndex(
	ctx context.Context,
	nameOrAddress string,
) (uint64, error) {
	app := common.HexToAddress(nameOrAddress)
	sel := table.InputIndex.
		SELECT(table.InputIndex.Index).
		FROM(table.InputIndex).
		WHERE(
			table.InputIndex.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var index uint64
	err := row.Scan(
		&index,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (r *postgresRepository) UpdateInputIndex(
	ctx context.Context,
	nameOrAddress string,
) error {
	app := common.HexToAddress(nameOrAddress)
	index, err := r.GetInputIndex(ctx, nameOrAddress)
	if err != nil {
		return err
	}
	nextIndex := index + 1

	inputIndexInsertStmt := table.InputIndex.INSERT(
		table.InputIndex.ApplicationAddress,
		table.InputIndex.Index,
	).VALUES(
		postgres.Bytea(app.Bytes()),
		nextIndex,
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	sqlStr, args := inputIndexInsertStmt.
		ON_CONFLICT(table.InputIndex.ApplicationAddress).
		DO_UPDATE(postgres.SET(
			table.InputIndex.Index.SET(postgres.RawInt(fmt.Sprintf("%d", nextIndex))),
		)).Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)

	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return nil
}

func (r *postgresRepository) GetLastProcessedEspressoBlock(
	ctx context.Context,
	nameOrAddress string,
) (uint64, error) {
	app := common.HexToAddress(nameOrAddress)
	sel := table.EspressoBlock.
		SELECT(table.EspressoBlock.LastProcessedEspressoBlock).
		FROM(table.EspressoBlock).
		WHERE(
			table.EspressoBlock.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)
	var lastProcessedEspressoBlock uint64
	err := row.Scan(
		&lastProcessedEspressoBlock,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return lastProcessedEspressoBlock, nil
}

func (r *postgresRepository) UpdateLastProcessedEspressoBlock(
	ctx context.Context,
	nameOrAddress string,
	lastProcessedEspressoBlock uint64,
) error {
	app := common.HexToAddress(nameOrAddress)
	insertStmt := table.EspressoBlock.INSERT(
		table.EspressoBlock.ApplicationAddress,
		table.EspressoBlock.LastProcessedEspressoBlock,
	).VALUES(
		postgres.Bytea(app.Bytes()),
		lastProcessedEspressoBlock,
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	sqlStr, args := insertStmt.
		ON_CONFLICT(table.EspressoBlock.ApplicationAddress).
		DO_UPDATE(postgres.SET(
			table.EspressoBlock.LastProcessedEspressoBlock.SET(postgres.RawFloat(fmt.Sprintf("%d", lastProcessedEspressoBlock))),
		)).Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)

	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return nil
}

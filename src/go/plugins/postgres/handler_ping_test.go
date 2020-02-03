/*
** Zabbix
** Copyright (C) 2001-2019 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

package postgres

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"zabbix.com/pkg/plugin"
)

var sharedConn *postgresConn
var fakeConn *postgresConn

func NewConnPool(t testing.TB) (*postgresConn, error) {

	if sharedConn == nil {
		newConn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5433/postgres")
		if err != nil {
			return nil, err
		}
		// FIXME: for different PG versions
		sharedConn = &postgresConn{postgresPool: newConn, lastTimeAccess: time.Now(), version: "100006"}
	}
	return sharedConn, nil
}

func TestPlugin_pingHandler(t *testing.T) {
	var pingOK int64
	pingOK = 1
	// create pool or aquare conn from old pool for test
	sharedPool, err := NewConnPool(t)
	if err != nil {
		t.Fatal(err)
	}

	impl.Configure(&plugin.GlobalOptions{}, nil)

	type args struct {
		conn   *postgresConn
		params []string
	}
	tests := []struct {
		name    string
		p       *Plugin
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			fmt.Sprintf("pingHandler should return %d if connection is ok", postgresPingOk),
			&impl,
			args{conn: sharedPool},
			pingOK,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.pingHandler(tt.args.conn, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.pingHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.pingHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

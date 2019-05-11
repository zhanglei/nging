/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package model

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type FrpServerAndGroup struct {
	*dbschema.FrpServer
	Group   *dbschema.FrpGroup
	Running bool
}

func NewFrpServer(ctx echo.Context) *FrpServer {
	return &FrpServer{
		FrpServer: &dbschema.FrpServer{},
		Base:      base.New(ctx),
	}
}

type FrpServer struct {
	*dbschema.FrpServer
	*base.Base
}

func (f *FrpServer) Exists(name string, excludeIds ...uint) (bool, error) {
	cond := db.Cond{`name`: name}
	if len(excludeIds) > 0 {
		cond[`id`] = db.NotEq(excludeIds[0])
	}
	n, e := f.Param().SetArgs(cond).Count()
	return n > 0, e
}
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

package manager

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/settings"
)

func configPost(c echo.Context, groups ...string) error {
	m := model.NewConfig(c)
	formValues := c.Forms()
	mapx := echo.NewMapx(formValues)
	var configList map[string]map[string]*dbschema.Config
	if len(groups) > 0 {
		configList = map[string]map[string]*dbschema.Config{}
		defaults := settings.ConfigDefaults()
		for _, group := range groups {
			conf, ok := defaults[group]
			if !ok {
				continue
			}
			configList[group] = conf
		}
	} else {
		configList = settings.ConfigDefaults()
	}
	for group, configs := range configList {
		_, err := m.ListByGroup(group)
		if err != nil {
			if err == db.ErrNoMoreRows {
				for _, cfg := range configs {
					_, err := cfg.Add()
					if err != nil {
						return err
					}
				}
				continue
			}
			return err
		}
		encoder := settings.GetEncoder(group)
		gm := mapx.Get(group)
		if gm == nil {
			continue
		}
		added := map[string]int{}
		for k, v := range m.Objects() {
			v.CPAFrom(m.Config)
			added[v.Key] = k
			setting := gm.Get(v.Key)
			if setting == nil {
				continue
			}
			disabled := setting.Value(`disabled`)
			_v := setting.Get(`value`)
			if _v == nil {
				continue
			}
			value, err := settings.EncodeConfigValue(_v, v, encoder)
			if err != nil {
				return err
			}
			var n int64
			condition := db.And(
				db.Cond{`key`: v.Key},
				db.Cond{`group`: group},
			)
			n, err = m.Count(nil, condition)
			if err != nil {
				return err
			}
			if n < 1 {
				err = settings.InsertBy(c, configs, v.Key, value, disabled)
				if err != nil {
					return err
				}
			}
			set := echo.H{}
			if value != m.Value {
				set[`value`] = value
			}
			if _v.IsMap() {
				if m.Type != `json` {
					set[`type`] = `json`
				}
			} else if _v.IsSlice() {
				if m.Type != `list` {
					set[`type`] = `list`
				}
			} else {
				cfg, ok := configs[v.Key]
				if ok && cfg != nil && m.Type != cfg.Type {
					set[`type`] = cfg.Type
				}
				//set[`type`] = `text`
			}
			if len(disabled) > 0 && m.Disabled != disabled {
				set[`disabled`] = disabled
			}
			if len(set) > 0 {
				err = m.SetFields(nil, set, condition)
				if err != nil {
					return err
				}
			}
		}
		err = settings.InsertMissing(c, gm, added, configs, encoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func configGet(c echo.Context, groups ...string) error {
	m := model.NewConfig(c)
	if len(groups) > 0 {
		for _, group := range groups {
			cfg, err := m.ListMapByGroup(group)
			if err != nil {
				return err
			}
			c.Set(group, cfg) //Stored.base.siteName
		}
		return nil
	}
	_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Group(`group`)
	}, 0, -1)
	if err != nil {
		return err
	}
	for _, setting := range m.Objects() {
		group := setting.Group
		cfg, err := m.ListMapByGroup(group)
		if err != nil {
			return err
		}
		c.Set(group, cfg) //Stored.base.siteName
	}
	return nil
}

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

package dashboard

import (
	"github.com/webx-top/echo"
)

func NewBlock(content func(echo.Context) error) *Block {
	return &Block{content: content}
}

type Block struct {
	Tmpl    string //模板文件
	Footer  string //末尾模版或JS代码
	content func(echo.Context) error
}

func (c *Block) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *Block) SetContentGenerator(content func(echo.Context) error) *Block {
	c.content = content
	return c
}

type Blocks []*Block

func (c Blocks) Ready(block echo.Context) error {
	for _, blk := range c {
		if blk != nil {
			if err := blk.Ready(block); err != nil {
				return err
			}
		}
	}
	return nil
}

var blocks = Blocks{}

func BlockRegister(block ...*Block) {
	blocks = append(blocks, block...)
}

//BlockRemove 删除元素
func BlockRemove(index int) {
	if index < 0 {
		blocks = blocks[0:0]
		return
	}
	size := len(blocks)
	if size > index {
		if size > index+1 {
			blocks = append(blocks[0:index], blocks[index+1:]...)
		} else {
			blocks = blocks[0:index]
		}
	}
}

//BlockSet 设置元素
func BlockSet(index int, list ...*Block) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		blocks = append(blocks, list...)
		return
	}
	size := len(blocks)
	if size > index {
		blocks[index] = list[0]
		if len(list) > 1 {
			BlockSet(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		blocks = append(blocks, nil)
	}
	blocks = append(blocks, list...)
}

func BlockAll() Blocks {
	return blocks
}

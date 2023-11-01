# mazegohan
This is a maze game using Golang.
maze + gohan(ご飯) = mazegohan
But this software does not include the gohan.

# About game
## 内容
ランダムに生成され，指数関数的に難易度が増加する迷路をクリアしてください！
## 基本操作
各種矢印キー←↓↑→を使用して移動できます．
## 色が着いたセル
緑　: プレイヤー
青　: ゴール🏆
赤　: 死ぬ☠
黄　: 数マスの間，パワーアップ💪（数マス移動する間は赤　を踏んでも死なない）
白　: 壁
## 仕様上のヒント
パワーアップが持続するマスのカウントは，黄色　のマスに乗る度にリセットされ，その度に新しい値がランダムで与えられます．（最小６マス）
つまり，黄色　のマスに複数回乗っても パワーアップが持続するマスのカウントは，増加（加算）しません．

# LICENSE and Special Thanks
This project is based on [JoelOtter/termloop](https://github.com/JoelOtter/termloop) by Joel Auterson & termloop authors.
And it is built on [nsf/termbox-go](https://github.com/nsf/termbox-go) by termbox-go authors.
In particular, I used the Termloop example [pyramid.go](https://github.com/JoelOtter/termloop/blob/master/_examples/pyramid.go) as a reference from Termloop.

These project's the components are provided under the terms of the MIT license.

About license of the Termloop is [here](https://github.com/JoelOtter/termloop/blob/master/LICENSE)
About license of the Termbox is [here](https://github.com/nsf/termbox-go/blob/master/LICENSE)

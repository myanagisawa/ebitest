package kitchen

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	// Game ...
	Game struct {
		bg           Scene
		currentScene Scene
		debugText    *DebugText
		myCoin       Coin
		coins        []Coin
		WindowSize   Size
	}

	// Size ...
	Size struct {
		Width  int
		Height int
	}
)

// NewGame ...
func NewGame(w, h int) (*Game, error) {

	backGround, _ := NewBackGround()
	debugText, _ := NewDebugText()
	g := &Game{
		bg:        backGround,
		debugText: debugText,
		WindowSize: Size{
			Width:  w,
			Height: h,
		},
	}

	// 初期化時のシーンを登録
	// sink, _ := NewSink(&g.WindowSize)
	// g.currentScene = sink
	s := NewSceneImpl()
	g.currentScene = s

	// Coin
	c, _ := NewMyCoin()
	g.myCoin = c

	c2, _ := NewCoin()
	g.coins = append(g.coins, c2)

	// l, _ := NewSpotlight(300.0, 300.0, 200.0, 1)
	// g.light = *l
	return g, nil
}

// Update ...
func (g *Game) Update(r *ebiten.Image) error {
	const d = 16

	sw, sh := r.Size()
	//dbg := fmt.Sprintf("screen size: %d, %d", sw, sh)

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		fmt.Println("Game::Up")
		sh += d
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		fmt.Println("Game::Down")
		if 16 < sh && d < sh {
			sh -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		fmt.Println("Game::Left")
		if 16 < sw && d < sw {
			sw -= d
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		fmt.Println("Game::Right")
		sw += d
	}
	ebiten.SetScreenSize(sw, sh)

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		fmt.Println("Game::C")
		for i := 0; i < 10; i++ {
			c, _ := NewDebris(0)
			// 生成オブジェクトの衝突判定
			col := false
			for _, coin := range g.coins {
				if CollisionCoin(c, coin) {
					col = true
					break
				}
			}
			if !col {
				g.coins = append(g.coins, c)
			}
		}
	}

	if err := g.bg.Update(); err != nil {
		return err
	}
	if err := g.currentScene.Update(); err != nil {
		return err
	}
	if err := g.debugText.Update(); err != nil {
		return err
	}
	if err := g.myCoin.Update(); err != nil {
		return err
	}
	for _, c := range g.coins {
		if err := c.Update(); err != nil {
			return err
		}
	}
	// コインの衝突判定
	for _, c := range g.coins {
		if CollisionCoin(g.myCoin, c) {
			g.myCoin.Collision(&c)
			c.Collision(&g.myCoin)
		}
		_ = Dot(g.myCoin, c)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.bg.Draw(r)
	g.currentScene.Draw(r)
	g.debugText.Draw(r)
	g.myCoin.Draw(r)
	for _, c := range g.coins {
		c.Draw(r)
	}

	//ebitenutil.DebugPrint(r, dbg)
	return nil
}

// CollisionCoin ...
func CollisionCoin(coin1, coin2 Coin) bool {
	c1, c2 := coin1.Circle(), coin2.Circle()
	// (xc1-xc2)^2 + (yc1-yc2)^2 ≦ (r1+r2)^2
	var dx, dy, dr float64 = float64(c1.x - c2.x), float64(c1.y - c2.y), float64(c1.r + c2.r)
	if (dx*dx + dy*dy) <= dr*dr {
		return true
	}
	return false
}

// Dot ...
func Dot(coin1, coin2 Coin) float64 {
	//	x1*x2 + y1*y2
	p := coin1.Circle().x*coin2.Circle().x + coin1.Circle().y*coin2.Circle().y
	// log.Printf("Dot=%f", p)
	return p
}

/*

///////////////////////////////////////////////////
// パーティクル衝突後速度位置算出関数
//   pColliPos_A : 衝突中のパーティクルAの中心位置
//   pVelo_A     : 衝突の瞬間のパーティクルAの速度
//   pColliPos_B : 衝突中のパーティクルBの中心位置
//   pVelo_B     : 衝突の瞬間のパーティクルBの速度
//   weight_A    : パーティクルAの質量
//   weight_B    : パーティクルBの質量
//   res_A       : パーティクルAの反発率
//   res_B       : パーティクルBの反発率
//   time        : 反射後の移動時間
//   pOut_pos_A  : パーティクルAの反射後位置
//   pOut_velo_A : パーティクルAの反射後速度ベクトル
//   pOut_pos_B  : パーティクルBの反射後位置
//   pOut_velo_B : パーティクルBの反射後速度ベクトル

bool CalcParticleColliAfterPos(
   D3DXVECTOR3 *pColliPos_A, D3DXVECTOR3 *pVelo_A,
   D3DXVECTOR3 *pColliPos_B, D3DXVECTOR3 *pVelo_B,
   FLOAT weight_A, FLOAT weight_B,
   FLOAT res_A, FLOAT res_B,
   FLOAT time,
   D3DXVECTOR3 *pOut_pos_A, D3DXVECTOR3 *pOut_velo_A,
   D3DXVECTOR3 *pOut_pos_B, D3DXVECTOR3 *pOut_velo_B
)
{
   FLOAT TotalWeight = weight_A + weight_B; // 質量の合計
   FLOAT RefRate = (1 + res_A*res_B); // 反発率
   D3DXVECTOR3 C = *pColliPos_B - *pColliPos_A; // 衝突軸ベクトル
   D3DXVec3Normalize(&C, &C);
   FLOAT Dot = D3DXVec3Dot( &(*pVelo_A-*pVelo_B), &C ); // 内積算出
   D3DXVECTOR3 ConstVec = RefRate*Dot/TotalWeight * C; // 定数ベクトル

   // 衝突後速度ベクトルの算出
   *pOut_velo_A = -weight_B * ConstVec + *pVelo_A;
   *pOut_velo_B = weight_A * ConstVec + *pVelo_B;

   // 衝突後位置の算出
   *pOut_pos_A = *pColliPos_A + time * (*pOut_velo_A);
   *pOut_pos_B = *pColliPos_B + time * (*pOut_velo_B);

   return true;
}
*/

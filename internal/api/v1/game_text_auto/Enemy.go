package game_text_auto

import (
	"github.com/kamva/mgm/v3"
	"github.com/mozillazg/go-pinyin"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"time"
)

// Enemy 敌人
// spiritAttack 神识：角色的灵魂力，用来计算角色的灵魂攻击力
//
// spiritDefence 神防：角色的神识防御力
//
// bleed 血气：角色生命值
//
// strong 力量：角色的力量，用来计算角色的物理攻击力
//
// shooting 命中率：角色的技巧，用来计算角色的命中率、必杀率和大部分技能的触发率
//
// attackSpeed 攻速：角色的速度。If的追击阈值是5，也就是说，当一名角色的速度高于敌方5点及时，该角色可以在敌方攻击后再攻击一次。速度也是影响角色回避率的属性
//
// dodge 闪避：角色的幸运，主要影响必杀回避（运%），对命中与回避也有些许影响（运/2%）
//
// defence 防御：角色的物理防御
//
// moveSpeed 移动力：角色与一回合内在平地可以移动的格子数量，基础上限为10（算上鞋子）
//
// level 星级
type Enemy struct {
	mgm.DefaultModel `bson:",inline"`
	NameZW             string `json:"name_zw"`
	NamePY             string `json:"name_py"`
	Description      string `json:"description"`
	SpiritAttack int64 `json:"spiritAttack"`
	SpiritDefence int64 `json:"spiritDefence"`
	Bleed int64 `json:"bleed"`
	Strong int64 `json:"strong"`
	Shooting int64 `json:"shooting"`
	AttackSpeed int64 `json:"attackSpeed"`
	Dodge int64 `json:"dodge"`
	Defence int64 `json:"defence"`
	MoveSpeed int64 `json:"moveSpeed"`
	Level int64 `json:"level"`
}

func NewEnemy(name, desc string) *Enemy {
	return &Enemy{
		NameZW:        name,
		NamePY:        EnemyNameZW2PY(name),
		Description:   desc,
		SpiritAttack:  1,
		SpiritDefence: 1,
		Bleed:         1000,
		Strong:        1,
		Shooting:      1,
		AttackSpeed:   1,
		Dodge:         1,
		Defence:       1,
		MoveSpeed:     1,
		Level: 1,
	}
}

func EnemyNameZW2PY(name string) string {
	namePy := ""
	pyStr := pinyin.NewArgs()
	for _, item := range pinyin.Pinyin(name, pyStr) {
		namePy += item[0]
	}
	return namePy
}

func RoundGenerateEnemy(baseProperty int64) (enemy Enemy, err error) {
	var enemies []Enemy
	err = mgm.Coll(&Enemy{}).SimpleFind(&enemies, bson.M{})
	if err != nil {
		return enemy, err
	}

	rand.Seed(time.Now().Unix())
	roundIndex := rand.Int63n(int64(len(enemies)))
	enemy = enemies[roundIndex]

	return Enemy{
		NameZW: enemy.NameZW,
		NamePY:        enemy.NamePY,
		Description: enemy.Description,
		SpiritAttack:  roundValue(10) + baseProperty,
		SpiritDefence: roundValue(10) + baseProperty,
		Bleed:         roundValue(10) + baseProperty,
		Strong:        roundValue(10) + baseProperty,
		Shooting:      roundValue(10) + baseProperty,
		AttackSpeed:   roundValue(10) + baseProperty,
		Dodge:         roundValue(10) + baseProperty,
		Defence:       roundValue(10) + baseProperty,
		MoveSpeed:     roundValue(10) + baseProperty,
		Level: roundValue(10),
	}, nil
}

func roundValue(max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max)
}
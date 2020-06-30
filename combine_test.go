package paodekuai

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
	// "games.agamestudio.com/jxjyhj/model"
)

func TestAlls(t *testing.T) {
	lastKind := CreateCardNo()
	lastKind.kind = KindOfCard_AAABB
	lastKind.max = 3
	lastKind.min = 3
	lastKind.num = 0
	lastKind.main = 1
	// v := selectOut([]int{12, 11, 10, 9, 0, 1, 13, 26}, nil, [][]int{}) //[32 45 19] 7 3 3 5 1
	// v := selectOut([]int{5, 38, 18, 10, 36, 3, 32, 26, 22, 40, 47, 41, 2, 49, 33, 34}, lastKind, [][]int{})
	// v := selectOut([]int{30, 21, 5, 18, 29, 6, 39, 41, 35, 47, 17, 34, 37, 45, 50, 48}, lastKind, [][]int{})
	// v := selectOut([]int{34, 19, 12, 13, 44, 29, 43, 26, 41, 22, 37, 10, 39, 20, 16, 7}, lastKind, [][]int{})
	// v := selectOut([]int{40, 12, 20, 17, 21, 18}, lastKind, [][]int{{5}, {6, 7}}) // 报单出最大
	// v := selectOut([]int{4, 0}, nil, [][]int{{5, 8}, {6, 7}})
	// v := selectOut([]int{46, 26, 13, 33, 47}, lastKind, [][]int{{5, 8}, {6, 7}})
	// v := selectOut([]int{20, 11, 10, 23, 33, 24, 46, 12}, lastKind, [][]int{{37}, {6, 7}})
	// v := selectOut([]int{20, 11, 21, 34, 33, 8, 46, 12}, lastKind, [][]int{{37}, {6, 7}})
	// v := selectOut([]int{17, 25}, nil, [][]int{{8, 50, 6, 7, 3, 47, 49, 32, 9, 33, 51, 4, 16, 22, 31},
	//	{43, 39, 35, 34, 19, 29, 45, 44, 15, 23, 10, 11, 30, 26, 46}})
	// v := selectOut([]int{6, 8, 11, 12, 13, 15, 19, 23, 24, 31, 32, 35, 37, 47, 48, 49}, lastKind,
	//	[][]int{{30, 7, 46, 5, 22, 20, 43, 17, 21, 9, 51},
	//		{33, 25, 26, 45, 36, 38, 10, 50, 18, 34, 40, 4, 44, 29, 39, 42}})
	// v := selectOut([]int{6, 8, 15, 19, 23, 31, 32, 35, 47, 48, 49}, nil,
	//	[][]int{{30, 7, 46, 5, 22, 20, 43, 17, 21, 9, 51},
	//		{33, 25, 26, 45, 36, 38, 10, 50, 18, 34, 40, 4, 44, 29, 39, 42}})
	// v := selectOut([]int{37, 29, 35, 46, 42, 11, 40, 12, 25, 51, 5, 16, 18, 17}, nil,
	//	[][]int{{34, 21, 26, 38, 13},
	//		{2, 43, 6, 36, 24, 45, 50, 33, 10, 4, 15, 47, 41, 9, 49}})
	// v := selectOut([]int{13, 17, 26, 39, 43}, nil,
	//	[][]int{{30, 38, 47, 32, 33, 18, 51, 10, 48, 29, 50}, 38, 43, 25, 4, 39
	//		{11, 16, 21, 37, 19, 7, 42, 20, 25, 22, 46}})
	// firstOut出大对问题，原因是原处理逻辑带牌已经固定一次出牌，存在多个三带情况下，排序会预先将大三带带上小对，排序到小三带情况下只能带大对。
	// v := selectOut([]int{2, 9, 10, 11, 12, 15, 23, 26, 28, 34, 35, 36, 39, 43, 48, 51}, nil,
	//	[][]int{{38, 43, 25, 4, 39},
	//		{5, 22, 13, 17, 30, 15, 35, 11, 31, 41, 42, 8, 16, 21, 26}}) // 2 15 28 26 39
	// v := selectOut([]int{44, 21, 28, 9, 51, 40, 34, 11, 39, 12, 29, 35, 43, 20, 46, 10}, nil, [][]int{})
	/*v := selectOut([]int{3, 29, 16, 40, 42, 51}, lastKind,
	[][]int{{45, 6, 7, 5, 18, 46, 38},
		{25, 12, 37}})*/
	// v := selectOut([]int{17, 2, 11, 51, 7, 39, 4, 33, 21, 48, 23, 31, 28, 35, 20, 24}, nil, [][]int{})
	// v := selectOut([]int{7, 51, 5, 39, 46, 44, 33}, nil,
	//	[][]int{{20, 50, 38, 41, 37, 34, 40, 29, 36, 4, 23, 43}, {5, 44}})
	/*v := selectOut([]int{33, 3, 7, 42, 20, 16}, lastKind,
	[][]int{{23, 18, 19, 50, 40, 30, 45, 4, 17, 32, 46}, {5, 44}})*/
	//v, err := selectOut([]int{33, 3, 7, 42, 20, 16}, nil,
	//	[][]int{{23, 18, 19, 50, 40, 30, 45, 4, 17, 32, 46}, {5, 44}})
	/*v, err := selectOut([]int{5, 10, 18, 26, 31, 44}, nil,
	[][]int{{49, 48, 35, 37, 50, 6, 32, 23, 34, 7, 19},
		{11, 45, 20, 36, 40, 39, 9, 24, 29, 13, 8, 38, 33, 3, 21, 17}})*/
	/*v, err := selectOut([]int{36, 11}, nil,
	[][]int{{45},
		{20, 32, 44, 48, 34, 19}})*/ // 全大加报单出最大特例
	/*v, err := selectOut([]int{13, 12, 11, 37, 51, 25, 26, 39}, lastKind,
	[][]int{{4, 24, 16, 9, 2, 47, 49, 45, 6, 40, 10, 34, 15, 3, 41, 8},
		{30, 7, 46, 17, 35, 23, 36, 44, 20, 32, 18, 22, 29, 38, 5, 50}}) // 全大加三带带牌不够情况*/ // 13, 15, 28, 29, 41, 42, 51, 31, 16
	/*v, err, bo := selectOut([]int{25, 46, 6, 15, 22, 18, 19, 36, 2, 24, 32, 28, 30, 20, 41, 12}, nil,
	[][]int{{41, 8, 15, 3, 16, 9, 4, 24, 2, 45, 6, 40, 10, 34, 47, 49},
		{38, 5, 32, 18, 22, 29, 50, 35, 23, 30, 7, 46, 17, 20, 36, 44}}) */
	v, err, bo := selectOut([]int{4, 17, 30, 43, 10, 23}, lastKind, [][]int{{0, 1},
		{2, 3}})
	if err != nil {
		cardKindParam := createDefaultCardParam()
		card, bo := AutoShowCard(nil, []int{13, 15, 28, 29, 41, 42, 51, 31, 16}, cardKindParam,
			[][]int{{41, 8, 15, 3, 16, 9, 4, 24, 2, 45, 6, 40, 10, 34, 47, 49},
				{38, 5, 32, 18, 22, 29, 50, 35, 23, 30, 7, 46, 17, 20, 36, 44}})
		fmt.Println("here", card, bo)
	}
	log.Println("####################", err, bo)
	fmt.Println(v.outPokers([]int{25, 46, 6, 15, 22, 18, 19, 36, 2, 24, 32, 28, 30, 20, 41, 12})) // outCard 51, 25, 26, 39, 13, 12, 11, 37
	v.Print(os.Stdout)

}

var (
	//黑桃-51,50,49,48,47,46,45,44,43,42,41,40,39
	//红桃-38,37,36,35,34,33,32,31,30,29,28,27,26
	//梅花-25,24,23,22,21,20,19,18,17,16,15,14,13
	//方片-12,11,10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0
	// 排序后的牌库
	library = []int{51, 50, 49, 48, 47, 46, 45, 44, 43, 42, 41, 40, 39,
		38, 37, 36, 35, 34, 33, 32, 31, 30, 29, 28, 27, 26,
		25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13,
		12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	}
)

func TestShowCard(t *testing.T) { // 牌库生成器
	var opCards [3][16]int
	winCards := make([][16]int, 0, 1000)
	loserCards := make([][][]int, 0, 1000)
	winpos := make([]int, 0, 1000)
	for m := 0; m != 10000; m++ {
		fmt.Printf("第%d轮游戏开始\n", m)
		p := NewPoker(0)
		cards := p.GetDifficultlyCards(false, 3)
		for i := 0; i != 3; i++ {
			copy(opCards[i][:], cards[i])
		}
		bankPos := 0
		for i := 0; i < Normal_HandCardNum; i++ {
			for j := 0; j < MaxNumOfPlayer; j++ {
				if cards[j][i] == Red_Peach {
					bankPos = j
					break
				}
			}
		}
		othersCard := make([][]int, 2)
		lastKind := CreateCardNo()

		var flag bool
		olp := bankPos
		for i := bankPos; ; i = (i + 1) % 3 {
			othersCard[0] = othersCard[0][:0]
			othersCard[1] = othersCard[1][:0]
			othersCard[0] = append(othersCard[0], cards[(i+1)%3]...)
			othersCard[1] = append(othersCard[1], cards[(i+1)%3]...)
			if i == olp { // 自动出牌
				lastKind = nil
			}
			v, err, bo := selectOut(cards[i], lastKind, othersCard)
			log.Printf("\n%d号玩家出牌####################%v %v %v", i, lastKind, bo, err)
			var outCard []int
			if err != nil {
				cardKindParam := createDefaultCardParam()
				outCard, bo = AutoShowCard(lastKind, cards[i], cardKindParam,
					othersCard)
				fmt.Println("AutoShowCard", outCard, bo)
				v, err, bo = selectOut(outCard, lastKind, othersCard)
				fmt.Println(v.outPokers(cards[i]), cards[i], bo)
				v.Print(os.Stdout)
			} else {
				outCard = v.outPokers(cards[i])
				fmt.Println(outCard) // outCard
				v.Print(os.Stdout)
			}
			if len(outCard) != 0 {
				if i == olp {
					flag = true
				}
				olp = i
				lastKind = CreateCardNo()
				if v != nil {
					lastKind.kind = v.hintType()
					lastKind.min = v.min
					lastKind.num = v.hintLen()
				}
			}
			if flag && lastKind != nil && (lastKind.kind == KindOfCard_AAA || lastKind.kind == KindOfCard_AAABBB ||
				lastKind.kind == KindOfCard_AAAB || lastKind.kind == KindOfCard_AAABBBCD) {
				if len(cards[i]) != len(outCard) {
					fmt.Println("errOut", outCard, cards[i])
					panic(errors.New("AAA"))
				}
			}
			flag = false
			for _, v := range outCard {
				for k, v1 := range cards[i] {
					if v1 == v {
						cards[i] = append(cards[i][:k], cards[i][k+1:]...) // del
					}
				}
			}
			if len(cards[i]) == 0 {
				fmt.Println("win", i, opCards[i])
				winCards = append(winCards, opCards[i])
				winpos = append(winpos, i)
				loCards := make([][]int, 2)

				ii := 0
				for j := 0; j != 3; j++ {
					if j == i {
						continue
					}
					loCards[ii] = make([]int, 16)
					copy(loCards[ii], opCards[j][:])
					ii++
					if ii == 2 {
						break
					}
					//loCards = append(loCards, opCards[j][:])
				}
				loserCards = append(loserCards, loCards)
				break
			}
		}
	}
	fmt.Printf("WinCard: \n")
	for k, v := range winCards {
		fmt.Printf("{")
		for i, v1 := range v {
			if i != 15 {
				fmt.Printf("%d, ", v1)
			} else {
				fmt.Printf("%d", v1)
			}
		}
		fmt.Printf("}, // %d\n", k)
	}
	for k, v := range loserCards {
		fmt.Printf("{")
		for _, v1 := range v {
			fmt.Printf("{")
			for i, v2 := range v1 {
				if i != 15 {
					fmt.Printf("%d, ", v2)
				} else {
					fmt.Printf("%d", v2)
				}

			}
			fmt.Printf("},")
		}
		fmt.Printf("}, // %d\n", k)
	}
	fmt.Printf("{\n")
	for k, v := range winpos {
		fmt.Printf("%d, ", v)
		fmt.Printf(" // %d\n", k)
	}
	fmt.Printf("}\n")

}

func randomShuffle(src []int) []int {
	dest := make([]int, len(src))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

type intSort []int

func (h intSort) Len() int {
	return len(h)
}
func (h intSort) Less(i, j int) bool {
	if h[i] > h[j] {
		return true
	}
	return false
}
func (h intSort) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// 包含关系,排序然后删除 sort
func delCards(ini []int, del []int) []int {
	sort.Sort(intSort(del))
	fmt.Println("sort:", ini)
	liNum := len(ini)
	delNum := len(del)
	ret := make([]int, 0, liNum-delNum)
	j := 0
	for _, v := range ini {
		if j == delNum {
			ret = append(ret, v)
			continue
		}
		if v == del[j] {
			j++
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

type pdkcardInfo struct {
	winCards   [][16]int
	loserCards [][][]int
	wins       []int
	gameNum    []int
}

var mo = NewMgoSession()

func TestTrainWinCard(t *testing.T) {
	fmt.Printf("############\n")
	fmt.Printf("")
	fmt.Printf("\n")
	var cards [3][]int
	var opCards [3][16]int
	var wins pdkcardInfo
	winCards, loserCards, winsPos := mo.GetLibrary() // 100
	hitNum := 0
	for i := 0; i != 3; i++ {
		cards[i] = make([]int, 0, 16)
	}
	for m := 0; m != 100; m++ {
		for i := 0; i != 3; i++ {
			cards[i] = cards[i][:0]
		}
		//del := make([]int, 16) // 深copy
		//copy(del, winCards[m][:])
		robot := winsPos[m]
		cards[robot] = append(cards[robot], winCards[m]...)
		//fmt.Println(del)
		//ret := delCards(library, del)
		//randCard := randomShuffle(ret)
		fmt.Println("sort:", cards[robot], winCards[m])
		//fmt.Println("ret", randCard)
		k := 0
		//for i := robot + 1; i != robot; i = (i + 1) % 3 {
		for i := 0; i != 3; i++ {
			if i == robot {
				continue
			}
			cards[i] = append(cards[i], loserCards[m][k]...)
			//cards[i] = append(cards[i], randCard[k:k+16]...)
			//k += 16
			k++
		}
		for i := 0; i != 3; i++ {
			copy(opCards[i][:], cards[i])
		}
		for i := 0; i != 3; i++ {
			fmt.Printf("[%d] 玩家手牌 %v\n", i, cards[i])
		}
		bankPos := 0
		for i := 0; i < Normal_HandCardNum; i++ {
			for j := 0; j < MaxNumOfPlayer; j++ {
				if cards[j][i] == Red_Peach {
					bankPos = j
					break
				}
			}
		}
		othersCard := make([][]int, 2)
		lastKind := CreateCardNo()
		lastKind = nil
		var flag bool
		olp := bankPos
		for i := bankPos; ; i = (i + 1) % 3 {
			othersCard[0] = othersCard[0][:0]
			othersCard[1] = othersCard[1][:0]
			othersCard[0] = append(othersCard[0], cards[(i+1)%3]...)
			othersCard[1] = append(othersCard[1], cards[(i+1)%3]...)
			if i == olp { // 自动出牌
				lastKind = nil
			}
			v, err, bo := selectOut(cards[i], lastKind, othersCard)
			log.Printf("\n%d号玩家出牌####################%v %v", i, lastKind, bo, err)
			var outCard []int
			if err != nil {
				cardKindParam := createDefaultCardParam()
				outCard, bo = AutoShowCard(lastKind, cards[i], cardKindParam,
					othersCard)
				fmt.Println("AutoShowCard", outCard, bo)
				v, err, bo = selectOut(outCard, lastKind, othersCard)
				fmt.Println(v.outPokers(cards[i]), cards[i], bo)
				v.Print(os.Stdout)
			} else {
				outCard = v.outPokers(cards[i])
				fmt.Println(outCard) // outCard
				v.Print(os.Stdout)
			}
			if len(outCard) != 0 {
				if i == olp {
					flag = true
				}
				olp = i
				lastKind = CreateCardNo()
				if v != nil {
					lastKind.kind = v.hintType()
					lastKind.min = v.min
					lastKind.num = v.hintLen()
				}
			}
			if flag && lastKind != nil && (lastKind.kind == KindOfCard_AAA || lastKind.kind == KindOfCard_AAABBB ||
				lastKind.kind == KindOfCard_AAAB || lastKind.kind == KindOfCard_AAABBBCD) {
				if len(cards[i]) != len(outCard) {
					fmt.Println("errOut", outCard, cards[i])
					panic(errors.New("AAA"))
				}
			}
			flag = false
			for _, v := range outCard {
				for k, v1 := range cards[i] {
					if v1 == v {
						cards[i] = append(cards[i][:k], cards[i][k+1:]...) // del
					}
				}
			}
			if len(cards[i]) == 0 {
				if i == robot {
					hitNum++
				} else {
					wins.winCards = append(wins.winCards, opCards[i])
					wins.wins = append(wins.wins, i)
					wins.gameNum = append(wins.gameNum, m)
					loCards := make([][]int, 2)

					ii := 0
					for j := 0; j != 3; j++ {
						if j == i {
							continue
						}
						loCards[ii] = make([]int, 16)
						copy(loCards[ii], opCards[j][:])
						ii++
						if ii == 2 {
							break
						}
					}
					wins.loserCards = append(wins.loserCards, loCards)
				}
				fmt.Println("win", i)
				break
			}
		}
	}
	for i, v := range wins.winCards {
		fmt.Printf("赢家:%d,第%d局\n", wins.wins[i], wins.gameNum[i])
		fmt.Printf("{")
		for i, v1 := range v {
			if i != 15 {
				fmt.Printf("%d, ", v1)
			} else {
				fmt.Printf("%d", v1)
			}
		}
		fmt.Printf("},\n")
		fmt.Printf("改变前牌型%v\n", winCards[wins.gameNum[i]])
		mo.UpData(wins.gameNum[i], v[:], wins.loserCards[i], wins.wins[i])
	}
	fmt.Printf("命中概率:%v \n", hitNum)
}

func TestUpdataWinCard(t *testing.T) {
	m := NewMgoSession()
	m.InitInset(WinCards, LoserCards, winPos)
	// 追加
	/*m.Inset([]int{17, 2, 32, 41, 16, 19, 47, 43, 15, 8, 20, 25, 44, 48, 23, 5},
	[][]int{{45, 10, 21, 46, 26, 34, 3, 31, 12, 37, 39, 9, 24, 49, 35, 36}, {6, 50, 51, 29, 13, 4, 30, 18, 40, 11, 42, 7, 22, 28, 38, 33}}, 0)*/
	wincards, loserCards, winpos := m.GetLibrary()
	for k, v := range wincards {
		fmt.Printf("{")
		for i, v1 := range v {
			if i != 15 {
				fmt.Printf("%d, ", v1)
			} else {
				fmt.Printf("%d", v1)
			}
		}
		fmt.Printf("}, // %d\n", k)
	}
	for k, v := range loserCards {
		fmt.Printf("{")
		for _, v1 := range v {
			fmt.Printf("{")
			for i, v2 := range v1 {
				if i != 15 {
					fmt.Printf("%d, ", v2)
				} else {
					fmt.Printf("%d", v2)
				}

			}
			fmt.Printf("},")
		}
		fmt.Printf("}, // %d\n", k)
	}
	fmt.Printf("{\n")
	for k, v := range winpos {
		fmt.Printf("%d, ", v)
		fmt.Printf(" // %d\n", k)
	}
	fmt.Printf("}\n")
}

func TestUpdataPdkWinLibrary(t *testing.T) {
	m := NewMgoSession()
	wincards, loserCards, _ := m.GetLibrary()
	model.PdkWinLibraryInit(wincards, loserCards)
}

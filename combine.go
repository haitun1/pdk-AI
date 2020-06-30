package paodekuai

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

/*
   尝试新的一套跑得快机器人出牌策略，要求易懂简洁，AI够智能
   简单描述:
   先将所有牌型(单牌除外)提取出来，压牌额外提取一手能压的手牌(必须出,非最优出牌能管就拆牌)，然后通过出牌->去重->出牌组成不同出牌方式。
   通过手数选出最优解，然后出牌。由于没混牌全组合效率也可以接受。后续可根据不同玩法设计通用方法。
*/

type line struct {
	min int
	num int
	len int
	tks []*line
}

func (l *line) incr(m map[int]int) {
	for i := l.min; i < l.min+l.len; i++ {
		m[i] += l.num
	}
}

func (l *line) decr(m map[int]int) {
	for i := l.min; i < l.min+l.len; i++ {
		m[i] -= l.num
	}
}

func (l *line) deepEqual(r *line) bool {
	// 地址相等必相同
	if l == r {
		return true
	}
	if l.min != r.min || l.num != r.num || l.len != r.len {
		return false
	}
	if len(l.tks) != len(r.tks) {
		return false
	}
	for i := 0; i < len(l.tks); i++ {
		if !l.tks[i].deepEqual(r.tks[i]) {
			return false
		}
	}
	return true
}

// 权重
func (l *line) Type() int {
	switch l.num {
	case 4:
		// 炸弹
		return 7
	case 3:
		if l.len >= 2 {
			// 飞机
			return 6
		} else if l.tks == nil && l.min == 11 {
			return 8 // AAA
		}
		// 三牌
		return 5
	case 2:
		if l.len >= 2 {
			// 对连
			return 3
		}
		// 对牌
		return 2
	case 1:
		if l.len >= 5 {
			// 单连
			return 4
		}
		// 单牌
		return 1
	}
	return 0
}

// 三牌带牌问题 只有最后一手可以出带单牌或者不带 所以判断带牌需要知道本身手牌
func (l *line) canTake(pL int, lastShowCard *KindOfCard) int {
	if l.tks != nil {
		return 0
	}
	/*if lastShowCard != nil && isTakeSingle(lastShowCard.GetKind()) { // 上家出牌为三带一或飞机带单
		return l.len
	} else*/
	if l.num == 3 {
		if pL == l.hintLen()+l.len {
			return l.len
		} else if pL == l.hintLen() {
			return 0
		}
		return l.len * 2
	}
	return 0
}

func (l *line) canOut(pL int) bool {
	if l.num == 3 {
		ll := l.hintLen()
		if ll != (l.num+2)*l.len { // 非带二情况
			if ll != pL && l.hintType() != KindOfCard_KingBomb {
				// fmt.Printf("HERE %v %v %v %v %v\n", l.num, l.min, l.len, ll, pL)
				return false
			}
		}
	}
	return true
}

func (l *line) isTake(b int) bool {
	for i := 0; i < l.len; i++ {
		if b == l.min+i {
			return false
		}
	}
	return true
}

// HintType .
func (l *line) HintType() int {
	return l.hintType()
}

func (l *line) hintType() int {
	switch l.num {
	case 4:
		if len(l.tks) == 0 {
			return KindOfCard_FourBomb
		}

		return KindOfCard_AAAABCD
	case 3:
		var n int
		if l.len != 1 { //
			n = 3
		}
		if len(l.tks) == 0 {
			if l.min == 11 {
				return KindOfCard_KingBomb
			}
			return KindOfCard_AAA + n
		}
		if len(l.tks) == 1 && l.tks[0].num == 1 {
			return KindOfCard_AAAB + n
		}
		return KindOfCard_AAABB + n
	case 2:
		if l.len == 1 {
			return KindOfCard_Double
		}
		return KindOfCard_AABBCC
	case 1:
		if l.len == 1 {
			return KindOfCard_Sigle
		}

		return KindOfCard_ABCDE
	}
	return KindOfCard_No
}

// HintLen .
func (l *line) HintLen() (n int) {
	return l.hintLen()
}

func (l *line) hintLen() (n int) {
	for _, v := range l.tks {
		n += v.hintLen()
	}
	n += int(l.num * l.len)
	return
}

// GetMin .
func (l *line) GetMin() int { return l.min }

// l > o 不同类型，小于等于为false 大于true
func (l *line) isBattle(o *line) bool {
	tp := l.hintType()
	ts := o.hintType()
	tpl := l.hintLen()
	tsl := o.hintLen()
	if tp == ts { // 肯定不是AAA
		if tpl == tsl {
			return l.min > o.min
		}
		return false
	}
	if tp == KindOfCard_KingBomb {
		return true
	}
	if tp == KindOfCard_AAA && (ts == KindOfCard_AAAB || ts == KindOfCard_AAABB) || // 模糊类型匹配
		(tp == KindOfCard_AAABBB && (ts == KindOfCard_AAABBBCD || ts == KindOfCard_AAABBBCCDD)) {
		return l.min > o.min
	}
	return false
}

func (l *line) isBigThan(lastShowCard *KindOfCard) bool {
	if lastShowCard == nil {
		return true
	}
	tp := l.hintType()
	if tp == lastShowCard.GetKind() {
		if l.hintLen() == lastShowCard.GetNum() { // 同类型，同数量
			return l.min > lastShowCard.GetMin()
		}
		return false
	}
	if tp == KindOfCard_KingBomb {
		return true
	}
	if tp == KindOfCard_FourBomb && lastShowCard.GetKind() != KindOfCard_FourBomb && lastShowCard.GetKind() != KindOfCard_KingBomb {
		return true
	}
	if tp != lastShowCard.GetKind() || l.hintLen() != lastShowCard.GetNum() {
		return false
	}
	return l.min > lastShowCard.GetMin()
}

func (l *line) outType(m int) int {
	if l.len > 1 || (l.len > 1 && l.num == 3) {
		if l.tks != nil && l.min <= m+0x8 {
			return 5 // 飞机带
		}
		if l.min <= m+0x8-l.num*2 {
			return 4 // 成连
		}
	}
	if l.tks != nil && l.min <= m {
		return 3 // 可带
	}
	for _, v := range l.tks {
		if v.outType(m) == 1 {
			return 2 // 被带
		}
	}
	if l.min <= m {
		return 1 // 普通
	}
	return 0
}

func (l *line) firstTake() *line {
	for _, v := range l.tks {
		return v // 最小的
	}
	return nil
}

func (l *line) Print(w io.Writer) {
	var trans = map[int]string{
		0:  "3",
		1:  "4",
		2:  "5",
		3:  "6",
		4:  "7",
		5:  "8",
		6:  "9",
		7:  "10",
		8:  "J",
		9:  "Q",
		10: "K",
		11: "A",
		12: "2",
	}

	l.nestPrint(w, trans, "")

	if l != nil {
		for _, v := range l.tks {
			v.nestPrint(w, trans, " 带")
		}
	}

	// fmt.Fprintln(w)
}

func (l *line) nestPrint(w io.Writer, trans map[int]string, take string) {
	if l != nil {
		switch l.Type() {
		case 6:
			fmt.Fprintf(w, "%d飞机: ", l.len)
		case 5:
			fmt.Fprint(w, "三牌: ")
		case 7:
			fmt.Fprint(w, "四牌: ")
		case 3:
			fmt.Fprintf(w, "%d对连: ", l.len)
		case 4:
			fmt.Fprintf(w, "%d单连: ", l.len)
		case 8:
			fmt.Fprint(w, "AAA炸弹: ")
		case 2:
			fmt.Fprintf(w, "%s对牌: ", take)
		case 1:
			fmt.Fprintf(w, "%s单牌: ", take)
		}

		for i := l.min; i < l.min+l.len; i++ {
			fmt.Fprint(w, trans[i], ", ")
		}
	} else {
		fmt.Fprintf(w, "没有找到该出手牌 ")
	}
}

// OutPokers .
func (l *line) OutPokers(poker []int) []int {
	return l.outPokers(poker)
}

func (l *line) outPokers(poker []int) []int {
	pok := make([]int, len(poker))
	copy(pok, poker)
	var ret []int
	if l != nil {
		for i := l.min; i != l.min+l.len; i++ {
			for j := 0; j != l.num; j++ {
				for k, v := range pok {
					if realVle(v) == i {
						ret = append(ret, v)
						pok = append(pok[:k], pok[k+1:]...) // del
						break
					}
				}
			}
		}
		for _, ll := range l.tks {
			for i := ll.min; i != ll.min+ll.len; i++ {
				for j := 0; j != ll.num; j++ {
					for k, v := range pok {
						if realVle(v) == i {
							ret = append(ret, v)
							pok = append(pok[:k], pok[k+1:]...) // del
							break
						}
					}
				}
			}
		}
	}
	return ret
}

type cutex struct {
	l []*line
	h int   // 手数
	b int   // 炸弹
	m []int // 哪几手非最大（方便测试人员，后在selectOut用于确认是否最大）
	f bool  // 四带标志（炸弹带）
}

func (c *cutex) takeOther(max int) bool {
	var n int
	for k := len(c.l) - 1; k >= 0; k-- {
		if c.l[k].len != 1 || c.l[k].num > 2 {
			break
		}
		if c.l[k].min >= max {
			continue
		}
		if n++; n >= 2 {
			return true
		}
	}
	return false
}

/*
   通过他人手牌推算出目前自己本手手牌是否为最大，最大不计入手数
*/
type enemy struct {
	otherPokers [][]int
}

func newEnemy(pokers [][]int) *enemy {
	var e enemy
	e.otherPokers = make([][]int, len(pokers))
	for i, v := range pokers {
		e.otherPokers[i] = append(e.otherPokers[i], v...)
	}
	return &e
}

func (e *enemy) isMax(ls *line) bool {
	for _, v := range e.otherPokers {
		m := ptm(v)
		lines := findBoom(m, 4)
		tp := ls.hintType()
		for _, v1 := range lines {
			if tp != KindOfCard_KingBomb && tp != KindOfCard_FourBomb {
				return false
			} else if v1.min > ls.min { // 不符合继续循环
				return false
			}
		}
		if len(v) < ls.hintLen() {
			continue
		}
		switch tp {
		case KindOfCard_ABCDE:
			lines = findPlane(m, 1, 5)
			for _, v1 := range lines {
				if v1.isBattle(ls) {
					return false
				}
			}
		case KindOfCard_AABBCC:
			lines = findPlane(m, 2, 2)
			for _, v1 := range lines {
				if v1.isBattle(ls) {
					return false
				}
			}
		case KindOfCard_AAABBB, KindOfCard_AAABBBCD, KindOfCard_AAABBBCCDD:
			lines = findPlane(m, 3, 2)
			for _, v1 := range lines {
				if v1.isBattle(ls) {
					return false
				}
			}
		default:
			lines = findSameType(m, ls.num)
			for _, v1 := range lines {
				if v1.isBattle(ls) {
					return false
				}
			}
		}

	}
	return true
}

type cutexSlice []*cutex

func (c cutexSlice) Len() int {
	return len(c)
}
func (c cutexSlice) Less(i, j int) bool {
	if c[i].h != c[j].h {
		return c[i].h < c[j].h // 少->多
	}
	if c[i].b != c[j].b {
		return c[i].b > c[j].b // 8888 7 6 全场最大时不带着出
	}
	return len(c[i].l) < len(c[j].l) // 888 7 6 全场最大时带着出
}
func (c cutexSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func realVle(v int) int { // 0-12 11:A 12:2 跑得快不包含大小王
	ret := v % PER_CARD_COLOR_MAX
	if ret < 2 {
		ret += PER_CARD_COLOR_MAX
	}
	return ret - 2 // 0-12
}

func ptm(in []int) map[int]int {
	m := make(map[int]int)
	for _, v := range in {
		m[realVle(v)]++
	}
	return m
}

func haveLine(m map[int]int, num int, len int, min int) bool {
	if min+len > 12 { // 2不能带
		return false
	}
	for i := min; i < min+len; i++ {
		if m[i] < num {
			return false
		}
	}
	return true
}

func makeLine(m, n, l int, s ...*line) *line {
	return &line{
		min: m, // 最小牌值
		num: n, // 单张数量
		len: l, // 长度
		tks: s,
	}
}

// 炸弹
func findBoom(m map[int]int, n int) (res []*line) {
	for i := 0; i <= 10; i++ {
		if m[i] >= n {
			res = append(res, makeLine(i, m[i], 1))
			res[len(res)-1].decr(m)
		}
	}

	if m[11] == 3 {
		res = append(res, makeLine(11, m[11], 1))
		res[len(res)-1].decr(m) // 砍掉AAA炸弹
	}

	return
}

// 单连
func findPlane(m map[int]int, num int, len int) (res []*line) {
	for i := len; ; i++ {
		var flag bool // 默认false 单联判断成功就继续
		for j := 0; j <= 11; j++ {
			if haveLine(m, num, i, j) {
				res = append(res, makeLine(j, num, i))
				flag = true
			}
		}
		if !flag {
			break
		}
	}
	return
}

func findCommon(m map[int]int, n int) (res []*line) {
	for i := 0; i <= 12; i++ {
		if m[i] >= n {
			res = append(res, makeLine(i, m[i], 1))
		}
	}
	return
}

func findSameType(m map[int]int, n int) (res []*line) {
	for i := 0; i <= 12; i++ {
		if m[i] >= n {
			res = append(res, makeLine(i, n, 1))
		}
	}
	return
}

// 遍历选定某连，找出与该连能同时存在的连，递归以上操作直至找不到，这些选定的连则是一种合法选择
// 请给参数t指向的切片足够的容量，如此以来在递归调用中将不会因扩展容量而重新申请内存
// 可以把递归想象成前序遍历树，由根向左，t用来记录由根到当前的所有节点
// 参数l和局部变量ll用来表示当前节点的子节点
// 参数m用来限制子节点的生成
// 参数d用来记录深度
func selectLine(m map[int]int, l []*line, t []*line, s *[][]*line, d int) {
	for _, v := range l {
		t = append(t[:d], v)

		v.decr(m)

		var ll []*line
		for _, vv := range l {
			// 不考虑自己（增加为压而拆后这里存在再找自己有效的情况，譬如对2被拆v是单2）
			if v == vv {
				continue
			}
			if haveLine(m, vv.num, vv.len, vv.min) {
				ll = append(ll, vv) // 能匹配的向来不多，这里不事先make
			}
		}

		if len(ll) == 0 {
			// 必须拷贝出来，因为参数t随着递归会变，前个合法选择会被下个覆盖
			tt := make([]*line, 0, len(t))
			tt = append(tt, t...)
			*s = append(*s, tt)
		} else {
			selectLine(m, ll, t, s, d+1)
		}

		v.incr(m)
	}
}

// 0： 相等， 1： 大于， -1： 小于
type compareFunc func(p, q *line) int

type multiSorter struct {
	ls []*line
	fs []compareFunc
}

func (ms *multiSorter) Sort(ls []*line) {
	ms.ls = ls
	sort.Sort(ms)
}

func (ms *multiSorter) Len() int {
	return len(ms.ls)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.ls[i], ms.ls[j] = ms.ls[j], ms.ls[i]
}

// 小 -> 大
func (ms *multiSorter) Less(i, j int) bool {
	for k := 0; k < len(ms.fs); k++ {
		if n := ms.fs[k](ms.ls[i], ms.ls[j]); n == 1 {
			return false
		} else if n == -1 {
			return true
		}
	}
	return false
}

func orderedBy(fs ...compareFunc) *multiSorter {
	return &multiSorter{
		fs: fs,
	}
}

// selectLine 后很容易[33 44][44 33]重复
// 但[33]均属于相同实例，排序后可以通过地址快速判定是否相同
// combineTake 后很容易带与被带均相同但由于带牌是拷贝而地址不等
// 必要时可为 selectLine 后单独实现 fastEqualLine 里面仅比较地址
func equalLine(l []*line, r []*line) bool {
	if len(l) != len(r) {
		return false
	}
	for i := 0; i < len(l); i++ {
		if !l[i].deepEqual(r[i]) {
			return false
		}
	}
	return true
}

func distinctLine(s [][]*line) [][]*line {
	// 排序便于比较相同
	for _, v := range s {
		orderedBy(rSortType, rSortMin, rSortLen).Sort(v)
	}
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if equalLine(s[i], s[j]) {
				s = append(s[:j], s[j+1:]...)
				j--
			}
		}
	}
	return s
}

func isTakeSingle(pokerType int) bool {
	if pokerType == KindOfCard_AAAB || pokerType == KindOfCard_AAABBBCD {
		return true
	}
	return false
}

func moreSplitTake(s *[][]*line, pL int, lastShowCard *KindOfCard) {
	ls := len(*s)

	for i := 0; i < ls; i++ {
		var ct int    // 可带牌数量
		p := (*s)[i]  // 始终指向即将拆对的
		pos := len(p) // 初始位置为长度预防无单对

		for k, v := range p {
			ct += v.canTake(pL, lastShowCard)

			if pos == len(p) && v.len == 1 && v.num <= 2 {
				pos = k // 记录首个单对
			}
		}

		// 单对不分类型按大小排序
		orderedBy(rSortMin).Sort(p[pos:])

		var bt int  // 被带牌数量
		var t []int // 记录对牌位置

		for j := len(p) - 1; j >= 0; j-- {
			if p[j].len != 1 || p[j].num > 2 {
				break
			}

			if bt++; bt >= ct {
				break
			}

			if p[j].num == 2 {
				t = append(t, j)
				bt++ // 对牌充当两单牌
			}
		}

		// 从小到大拆对
		for _, v := range t {
			n := make([]*line, 0, len(p)+1)
			l := makeLine(p[v].min, 1, 1)
			n = append(n, p[:v]...)
			n = append(n, p[v+1:]...)
			n = append(n, l, l)

			// 新增一种切牌
			*s = append(*s, n)
			orderedBy(rSortType, rSortMin, rSortLen).Sort(p[pos:])
			p = n // 下次以本次结果为基础
		}
		orderedBy(rSortType, rSortMin, rSortLen).Sort(p[pos:])
	}
}

func combineTake(s *[][]*line, n int, pL int, lastShowCard *KindOfCard) {
	ln := len(*s)

	for i := n; i < ln; i++ {
		for j := 0; j < len((*s)[i]); j++ {
			// 已没有可带牌
			if (*s)[i][j].num != 3 || (*s)[i][j].min == 11 { // 不考虑拆开四炸,AAA
				break
			}
			// 此牌已带过
			ct := (*s)[i][j].canTake(pL, lastShowCard)
			if ct == 0 {
				continue
			}

			// 记录被带单或对的位置
			is := make([]int, 0, ct)
			var id []int
			var cd int
			if ct%2 == 0 && ct != 0 {
				cd = ct / 2
				id = make([]int, 0, cd)
			}
			for k := len((*s)[i]) - 1; k >= 0; k-- {
				// 已没有单或对
				if (*s)[i][k].len != 1 || (*s)[i][k].num > 2 {
					break
				}
				// 校验带牌与被带牌不能相同
				if !(*s)[i][j].isTake((*s)[i][k].min) {
					continue
				}
				if (*s)[i][k].num == 1 && len(is) < int(ct) {
					is = append(is, k)
				}
				if (*s)[i][k].num == 2 && len(id) < int(cd) {
					id = append(id, k)
				}
				// 已找够单和对
				if len(is) == int(ct) && len(id) == int(cd) {
					break
				}
			}

			if len(is) == int(ct) {
				*s = append(*s, pickupTake((*s)[i], j, is))
			}
			if cd != 0 && len(id) == int(cd) {
				*s = append(*s, pickupTake((*s)[i], j, id))
			}
		}
	}

	// 有新增则递归
	if len(*s) > ln {
		combineTake(s, ln, pL, lastShowCard)
	}
}

// 带与被带之外为复制需考虑四点
// [#1][带][#2][被带2][#3][被带1][#4]
// 为便于理解刻意举需要考虑#3的例子
// [KKKAAA][555666][9][8][6][3]
// 当切成这样并考虑为56带38时
// m = 1; t = [53]
func pickupTake(s []*line, m int, t []int) []*line {
	r := make([]*line, 0, len(s)-len(t))
	r = append(r, s[:m]...) // #1
	tks := make([]*line, 0, len(t))
	for _, v := range t {
		tks = append(tks, s[v])
	}
	r = append(r, makeLine(s[m].min, s[m].num, s[m].len, tks...))
	for i := len(t) - 1; i >= 0; i-- {
		if i == len(t)-1 {
			r = append(r, s[m+1:t[i]]...) // #2
		} else {
			r = append(r, s[t[i+1]+1:t[i]]...) // #3
		}
	}
	return append(r, s[t[0]+1:]...) // #4
}

func splitThree(lines []*line) bool {
	min := -2
	flag := false
	for _, v := range lines {
		if v.num != 3 || v.min == 11 {
			flag = false
			break
		} else {
			if min == -2 {
				min = v.min
			} else if !flag && v.min-1 != min && v.min+1 != min { // 不是飞机的三带
				flag = true
			}
		}
	}
	return flag
}

// cutPokers pokers 自己手牌 lastShowCard 上家出牌信息 otherPokers 别人剩余手牌
func cutPokers(pokers []int, lastShowCard *KindOfCard, otherPokers [][]int) []*cutex {
	m := ptm(pokers)
	pLen := len(pokers)

	var lines []*line

	// line 代表该名玩家能组成的所有牌型
	// 考虑炸弹不可拆，可以先将炸弹剃掉，单独组成牌型
	lines = append(lines, findBoom(m, 4)...)             // 因为炸弹算分所以单独区分，不做拆开处理
	lines = append(lines, findPlane(m, 1, 5)...)         // 单连， 未去重
	lines = append(lines, findPlane(m, 2, 2)...)         // 对连
	lines = append(lines, findPlane(m, 3, 2)...)         // 飞机
	lines = append(lines, findCommon(m, 2)...)           // 四牌 三牌 对牌
	lines = append(lines, findPress(m, lastShowCard)...) // 如果上家有牌就拆开

	if splitThree(lines) { // 考虑极端情况，仅剩三带牌型
		lines = append(lines, findSameType(m, 2)...) // 将三带拆开
	}

	var s [][]*line
	t := make([]*line, 0, len(lines))

	selectLine(m, lines, t, &s, 0)

	// 去重
	s = distinctLine(s)
	for i := 0; i < len(s); i++ {
		n := len(s[i])
		for j := 0; j < n; j++ {
			s[i][j].decr(m)
		}
		// 拆四余三、拆三余对、拆对余单、原本的单牌
		s[i] = append(s[i], findCommon(m, 1)...)
		for j := 0; j < n; j++ {
			s[i][j].incr(m)
		}
		orderedBy(rSortType, rSortMin, rSortLen).Sort(s[i])
	}

	// 只有散牌
	if len(s) == 0 {
		s = append(s, findCommon(m, 1))
	} else {
		// 三带一要考虑可能拆对牌 跑得快可以三带两单
		//if isTakeSingle(lastShowCard.GetKind()) {
		moreSplitTake(&s, pLen, lastShowCard)
		//}

		// 带牌
		combineTake(&s, 0, pLen, lastShowCard)

		s = distinctLine(s)
	}

	cs := make([]*cutex, 0, len(s))
	enemy := newEnemy(otherPokers)
	var firstOutFlag bool
	if lastShowCard != nil || len(otherPokers) == 0 { // 由于首出牌敌人手牌传入为空，不限制会出错
		firstOutFlag = true
	}
	for _, v := range s {
		c := &cutex{
			l: v,
			h: len(v),
		}
		for k, vv := range v {
			if enemy.isMax(vv) && !firstOutFlag {
				c.h--
				/*} else if lastShowCard != nil && vv.isBigThan(lastShowCard) { // 被动出牌检测手牌是否都大于前着手牌 但检测实际意义不大,舍弃
				c.h--*/
			} else {
				c.m = append(c.m, k)
			}

			tp := vv.hintType()
			if tp == KindOfCard_FourBomb || tp == KindOfCard_KingBomb {
				c.b++
			}
			if !c.f && tp == KindOfCard_AAAABCD {
				c.f = true
			}
		}

		cs = append(cs, c)
	}
	sort.Sort(cutexSlice(cs))

	return cs
}

// SelectOut 对外开放接口
func SelectOut(pokers []int, lastShowCard *KindOfCard, otherPokers [][]int) (ls *line, retErr error, onlyRet bool) {
	return selectOut(pokers, lastShowCard, otherPokers)
}

func selectOut(pokers []int, lastShowCard *KindOfCard, otherPokers [][]int) (ls *line, retErr error, onlyRet bool) { // 选择手牌
	defer func() { // 增加容错处理
		if err := recover(); err != nil {
			retErr = errors.New("selectOut is err")
			ls = nil
		}
	}()
	if lastShowCard != nil && lastShowCard.num == 0 {
		//logger.Logger.Tracef("lastShowCard num is zero %v", lastShowCard)
		lastShowCard = nil
	}
	cs := cutPokers(pokers, lastShowCard, otherPokers)
	/*for _, v := range cs {
		for _, v1 := range v.l {
			fmt.Printf("CS: min:%v num:%v len:%v", v1.min, v1.num, v1.len)
			for _, v2 := range v1.tks {
				fmt.Printf(" tks:%v ", *v2)
			}
		}
		fmt.Printf(" cs h:%v b:%v m:%v f:%v \n", v.h, v.b, v.m, v.f)
	}*/
	pL := len(pokers)
	// var one *line
	var offset int
	// h==0 全是最大，找到直接出
	// h==1 先找最大，并记录首个非最大，若没有最大则出非最大，此时会赢
	for _, c := range cs {
		if c.h > 1 { // 不能一手出完就跳过
			break
		}
		for k, l := range c.l {
			if !l.isBigThan(lastShowCard) {
				continue
			}
			// 没有非最大或这个不是非最大，说明这个是最大 考虑如果剩下一手情况下
			if (len(c.m) == 0 || c.m[0] != k) && l.canOut(pL) { // 后续考虑代码优化
				/*if (l.hintType() == KindOfCard_KingBomb || l.hintType() == KindOfCard_FourBomb) && len(c.l) > 2 { //真实手数超过三手
					for i, v := range c.l {
						if len(c.m) == 0 {
							if i == k {
								continue
							}
						} else {
							if i == k || i == c.m[0] {
								continue
							}
						}
						ls = v
						if ls.canOut(pL) {
							ls = calcOut(ls, pokers, lastShowCard, otherPokers)
							fmt.Println("here")
							return
						}
					}
				}*/
				ls = l
				if ls.canOut(pL) {
					ls, onlyRet = calcOut(ls, pokers, lastShowCard, otherPokers, onlyRet)
					if ls.HintType() == KindOfCard_KingBomb && pL > 4 {
						continue
					} else if ls.HintType() == KindOfCard_FourBomb && pL > 3 {
						continue
					}
					return
				}
			}
			// 记录非最大
			if ls == nil {
				ls = l
			}
		}
		offset++
	}

	// 出非最大
	if ls != nil && ls.canOut(pL) {
		ls, onlyRet = calcOut(ls, pokers, lastShowCard, otherPokers, onlyRet)
		return //one, ret
	}

	if len(cs[offset:]) != 0 {
		cs = cs[offset:]
	}

	// 首出
	// var ls *line
	if lastShowCard == nil {

		ls, onlyRet, retErr = headOut(cs, pokers)
	} else {
		ls, onlyRet = pressOut(cs, lastShowCard)
	}
	ls, onlyRet = calcOut(ls, pokers, lastShowCard, otherPokers, onlyRet)
	return //ls, ret //pressOut(cs, lastShowCard)
}

// 这里主要限制AAA带牌问题
func containBoom(ls []*line) bool {
	if len(ls) == 0 {
		return false
	}
	orderedBy(sortMin, sortTks).Sort(ls)
	if ls[0].num == 3 && ls[0].min == 13 {
		return true
	}
	return false
}

func pressOut(cs []*cutex, lastShowCard *KindOfCard) (*line, bool) {
	var h int
	var ls []*line
	onlyRet := true

	// 首个能压手数中可出牌
	// 本身就是根据h排序过的
	for _, c := range cs {
		if h != 0 && c.h > h {
			/*if containBoom(ls) { // 判断管牌有没有包含
				ls = ls[:0]
			} else {*/
			break // continue
			//}
		}
		for _, l := range c.l {

			if !l.isBigThan(lastShowCard) {
				// fmt.Printf("cs: l%v  type:%v\n", *l, l.hintType())
				continue
			}

			ls = append(ls, l)
			h = c.h
		}
	}

	// 无牌可出
	if len(ls) == 0 {
		return nil, true
	}

	// 小 -> 小带
	orderedBy(sortMin, sortTks).Sort(ls)
	if len(ls) > 1 {
		onlyRet = false
	}
	return ls[0], onlyRet
}

func calcOut(ls *line, pokers []int, lastShowCard *KindOfCard, otherPokers [][]int, ret bool) (*line, bool) {
	if (len(otherPokers) != 0 && len(otherPokers[0]) == 1) && ls != nil && ls.hintType() == KindOfCard_Sigle { // 报单出最大
		m := ptm(pokers)
		if lastShowCard == nil { // 自动出牌出对牌
			for i := 0; i <= 11; i++ {
				if m[i] == 2 {
					return makeLine(i, m[i], 1), true
				}
			}
		}
		// 没有对牌，出大单
		if m[12] != 0 {
			return makeLine(12, m[12], 1), true
		}
		if m[11] != 0 {
			if m[11] == 3 {
				return makeLine(11, m[11], 1), true // AAA炸弹直接出
			}
			return makeLine(11, 1, 1), true

		}
		for i := 10; i >= 0; i-- {
			if m[i] != 0 {
				if m[i] == 4 {
					return makeLine(i, m[i], 1), true // 炸弹直接出
				}
				return makeLine(i, 1, 1), true
			}
		}
	}
	return ls, ret
}

func headOut(cs []*cutex, pokers []int) (l *line, onlyRet bool, retErr error) {
	defer func() { // 增加容错处理
		if err := recover(); err != nil {
			retErr = errors.New("headOut is err")
			l = nil
		}
	}()
	tps := make(map[int][]*line)

	// 最小逻辑牌值
	min := 12
	for _, v := range pokers {
		if n := realVle(v); n < min {
			min = n
		}
	}

	// 从最少手数中选出携带最小牌值的出牌按牌型记录
	for _, c := range cs {
		if c.h > cs[0].h {
			break
		}
		for _, l := range c.l {
			t := l.outType(min) // 暂不修改
			if t == 0 {         // 不携带最小牌值
				continue
			}
			if t == 2 && c.takeOther(l.min) {
				// 最小牌属于被带且当前切牌有别的被带可选
				tps[t] = append(tps[t], l.firstTake())
			} else {
				tps[t] = append(tps[t], l)
			}
		}
	}

	var i int
	var ls []*line

	// 飞机带 -> 成连 -> 可带 -> 被带 -> 普通
	for i = 5; i > 0; i-- {
		if v, ok := tps[i]; ok {
			ls = v
			break
		}
	}

	lLen := len(ls)
	pL := len(pokers)
	if lLen == 0 { // 没找到手牌,直接从cs拿牌出
		if ll := findCsPoker(cs, pL); ll != nil {
			l = ll
			return
		}
		l = nil
		retErr = errors.New("headOut is err")
		return
	}

	switch i {
	case 5:
		// 小 -> 长 -> 小带
		orderedBy(sortMin, rSortLen, sortTks).Sort(ls)
	case 4:
		// 小 -> 宽 -> 短
		orderedBy(sortMin, rSortNum, sortLen).Sort(ls)
	case 3:
		// 小带
		orderedBy(sortTks).Sort(ls)
	case 2:
		// 带 -> 小 -> 长
		orderedBy(rSortTake, sortMin, rSortLen).Sort(ls)
	case 1:
		// 宽
		orderedBy(rSortNum).Sort(ls)
	}
	lsLen := 0
	for _, v := range ls {
		if v.canOut(pL) {
			lsLen++
		}
	}
	if lsLen > 1 {
		onlyRet = true
	}
	firstOut := true
	for _, v := range ls {
		if v.canOut(pL) { // 主要限制是非出最后一手手牌外的禁止牌型
			l = v
			if firstOut && lsLen > 1 { // 不首出炸弹
				if v.HintType() == KindOfCard_KingBomb || v.HintType() == KindOfCard_FourBomb {
					continue
				}
				firstOut = false
			}
			return
		}
	}

	if ll := findCsPoker(cs, pL); ll != nil {
		l = ll
		return
	}
	// 走到这里已经出错了 1:ls==nil 2:ls出牌不符合规则

	l = nil
	retErr = errors.New("headOut is err") // 返回错误，使用老处理逻辑
	return                                // ls[10], ret
}

func findCsPoker(cs []*cutex, pl int) *line {
	// fmt.Printf("HERE\t")
	for _, v := range cs {
		for _, v1 := range v.l {
			if v1.canOut(pl) {
				// v1.Print(os.Stdout)
				return v1
			}
		}
	}

	return nil
}

func findPress(m map[int]int, lastShowCard *KindOfCard) []*line {
	if lastShowCard == nil {
		return nil
	}
	n := 4
	switch lastShowCard.GetKind() {
	case KindOfCard_Sigle:
		n = 1
	case KindOfCard_Double:
		n = 2
	case KindOfCard_AAA, KindOfCard_AAAB, KindOfCard_AAABB:
		n = 3
	}

	// 由于必须管的原因，必须将能管的拆出来备用。
	var ret []*line
	for i := 0; i <= 12; i++ {
		if m[i] >= n && i > lastShowCard.GetMin() {
			ret = append(ret, makeLine(i, n, 1))
			break // 效率问题 本身全拆后进行全组合导致性能低下
		}
	}

	return ret
}

var rSortLen = func(p, q *line) int {
	return sortLen(q, p)
}
var sortLen = func(p, q *line) int {
	if p.len == q.len {
		return 0
	}
	if p.len > q.len {
		return 1
	}
	return -1
}

var rSortNum = func(p, q *line) int {
	return sortNum(q, p)
}
var sortNum = func(p, q *line) int {
	if p.num == q.num {
		return 0
	}
	if p.num > q.num {
		return 1
	}
	return -1
}

var rSortMin = func(p, q *line) int {
	return sortMin(q, p)
}
var sortMin = func(p, q *line) int {
	if p.min == q.min {
		return 0
	}
	if p.min > q.min {
		return 1
	}
	return -1
}

var rSortTks = func(p, q *line) int {
	return sortTks(q, p)
}
var sortTks = func(p, q *line) int {
	/*defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err, len(p.tks), len(q.tks))
			p.Print(os.Stdout)
			q.Print(os.Stdout)
		}
	}()*/
	pL := len(p.tks)
	qL := len(q.tks)
	if pL != qL { // 带对牌跟两张单牌。只需考虑对牌需不需要拆开
		if pL == 0 {
			return 1
		} else if qL == 0 {
			return -1
		}
		if p.tks[0].min < q.tks[0].min {
			return -1
		}
		if p.tks[0].min > q.tks[0].min {
			return 1
		}
		if pL > qL {
			return 1
		}
		return -1
	}
	for i := 0; i < pL; i++ {
		if p.tks[i].min < q.tks[i].min {
			return -1
		}
		if p.tks[i].min > q.tks[i].min {
			return 1
		}
	}
	return 0
}

var rSortType = func(p, q *line) int {
	return sortType(q, p)
}
var sortType = func(p, q *line) int {
	if tp, tq := p.Type(), q.Type(); tp == tq {
		return 0
	} else if tp > tq {
		return 1
	}
	return -1
}

var rSortTake = func(p, q *line) int {
	return sortTake(q, p)
}
var sortTake = func(p, q *line) int {
	if p.tks != nil {
		if q.tks == nil {
			return 1
		}
	} else if q.tks != nil {
		return -1
	}
	return 0
}

// Imitate 模拟一局游戏并产生赢家
func Imitate(id int) (winCards []int, loseCards [][]int, winPos int) {
	var opCards [3][16]int
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
		v, err, _ := selectOut(cards[i], lastKind, othersCard)
		var outCard []int
		if err != nil {
			cardKindParam := createDefaultCardParam()
			outCard, _ = AutoShowCard(lastKind, cards[i], cardKindParam,
				othersCard)
			v, err, _ = selectOut(outCard, lastKind, othersCard)
			v.Print(os.Stdout)
		} else {
			outCard = v.outPokers(cards[i])
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

			return opCards[i][:], loCards, i
			break
		}
	}
	return
}

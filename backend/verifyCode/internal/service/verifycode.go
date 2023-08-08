package service

import (
	"context"
	"math/rand"
	pb "verifyCode/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode(int(req.Length), req.Type),
	}, nil
}

// 开放的被调用方法，用于区分类型
func RandCode(l int, t pb.TYPE) string {
	switch t {
	case pb.TYPE_DEFAULT:
		fallthrough
	case pb.TYPE_DIGIT:
		chars := "0123456789"
		return randCode(chars, l, 4)
	case pb.TYPE_LETTER:
		chars := "abcdefghijklmnopqrstuvwxyz"
		return randCode(chars, l, 5)
	case pb.TYPE_MIXED:
		chars := "0123456789abcdefghijklmnopqrstuvwxyz"
		return randCode(chars, l, 6)
	default:

	}
	return ""
}

// 随机数核心方法
// 一次随机多个随机位，分部分多次使用
func randCode(chars string, l, idxBits int) string {
	//计算有效的二进制数位 基于chars长度
	//推荐写死，因为chars固定，对应的位数长度也固定
	//idxBits = len(fmt.Sprintf("%b", len(chars)))
	//idxBits 验证码需要多少位
	//形成掩码 mask
	//例如，使用第六位   ：0000000000000111111
	idxMask := 1<<idxBits - 1 // 00000111111
	//63位可用多少次
	idxMax := 63 / idxBits
	//结果
	result := make([]byte, l)
	//生成随机字符
	//cache 随机位缓存  remain  当前还可以用几次
	for i, cache, remain := 0, rand.Int63(), idxMax; i < l; {
		if 0 == remain {
			cache, remain = rand.Int63(), idxMax
		}
		//利用掩码获取有效部位的随机数位
		if randIndex := int(cache & int64(idxMask)); randIndex < len(chars) {
			result[i] = chars[randIndex]
			i++
		}
		//使用下一组随机位
		cache >>= idxBits
		//减少一次使用次数
		remain--
	}

	return string(result)
}

// 随机的核心方法
//func randCode(chars string, l int) string {
//
//	charsLen := len(chars)
//	result := make([]byte, l)
//	for i := 0; i < l; i++ {
//		//生成[0,n)的整形随机数
//		randIndex := rand.Intn(charsLen)
//		result[i] = chars[randIndex]
//	}
//	rand.Intn(len(chars))
//	return string(result)
//}

package response

type ErrCode int

//go:generate stringer -type=ErrCode -linecomment

const (
	Ok               ErrCode = 0      // ok
	QueryParamsError ErrCode = 100002 // 请求参数错误，请检查后重试！
)

// 登录
const (
	ValidateCaptchaFail ErrCode = 100001 // 验证验证码错误
	SendCaptchaFail     ErrCode = 100002 // 发送验证码错误
	CaptchaRepeat       ErrCode = 100003 // 验证码已发送，请两分钟后重新发送
	CaptchaExpired      ErrCode = 100004 // 验证码已过期，请重新发送
	LoginFail           ErrCode = 100005 // 登录失败
)

// 视频下载
const (
	AcquiredVideoList ErrCode = 200001 // 获取视频列表数据失败
	DownloadError     ErrCode = 200002 // 下载失败
)

// 必应

const (
	AcquiredImagesFailed ErrCode = 300001 // 获取必应图片失败
)

// 网易云转换
const (
	ParseMetaInfoFail   ErrCode = 400000 // 解析音乐信息失败
	ProcessFail         ErrCode = 400001 // 网易云转换失败
	ChooseDirectoryFail ErrCode = 400002 // 选择目录失败
	ChooseFileFail      ErrCode = 400003 // 选择文件失败
	FolderNotFile       ErrCode = 400004 // 目录下面没有需要转换的文件
	ChooseFolderFail    ErrCode = 400005 // 选择文件夹失败
	ChooseFile          ErrCode = 400006 // 请选择需要转换的文件
	MusicNotExists      ErrCode = 400007 // 音乐文件不存在，请检查后重试
	MusicIncomplete     ErrCode = 400008 // 音乐文件未处理完成，请检查后重试
)

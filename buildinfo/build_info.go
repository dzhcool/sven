package buildinfo

/**
 * 需编译时候传递参数过来
 * 具体用法参考readme
 * @auther davis
 * @date 20210918
 */

var (
	build_time       string
	build_version    string
	build_go_version string
	build_author     string
)

func GetBuildTime() string {
	return build_time
}

func GetBuildVersion() string {
	return build_version
}

func GetBuildGoVersion() string {
	return build_go_version
}

func GetBuildAuthor() string {
	return build_author
}

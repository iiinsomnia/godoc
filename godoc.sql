/*
 Navicat Premium Data Transfer

 Source Server         : MySQL
 Source Server Type    : MySQL
 Source Server Version : 50718
 Source Host           : localhost:3306
 Source Schema         : godoc

 Target Server Type    : MySQL
 Target Server Version : 50718
 File Encoding         : 65001

 Date: 01/08/2017 21:56:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for go_category
-- ----------------------------
DROP TABLE IF EXISTS `go_category`;
CREATE TABLE `go_category` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) NOT NULL COMMENT '名称',
  `created_at` datetime NOT NULL COMMENT '添加时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of go_category
-- ----------------------------
BEGIN;
INSERT INTO `go_category` VALUES (1, '技术文档', '2017-07-30 18:35:55', '2017-07-30 18:35:55');
INSERT INTO `go_category` VALUES (2, '需求文档', '2017-07-30 18:36:05', '2017-07-30 18:36:05');
INSERT INTO `go_category` VALUES (3, 'API文档', '2017-07-30 18:36:13', '2017-07-30 18:36:13');
COMMIT;

-- ----------------------------
-- Table structure for go_doc
-- ----------------------------
DROP TABLE IF EXISTS `go_doc`;
CREATE TABLE `go_doc` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(50) NOT NULL COMMENT '名称',
  `category_id` int(11) NOT NULL COMMENT '类别ID',
  `project_id` int(11) NOT NULL COMMENT '项目ID',
  `label` varchar(50) NOT NULL COMMENT '标签',
  `markdown` mediumtext COMMENT 'markdown',
  `created_at` datetime NOT NULL COMMENT '添加时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_project` (`project_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='文档表';

-- ----------------------------
-- Records of go_doc
-- ----------------------------
BEGIN;
INSERT INTO `go_doc` VALUES (1, 'nginx 配置 location 及 rewrite 规则', 1, 1, 'nginx,location,rewrite', '# nginx 配置 location 及 rewrite 规则\r\n\r\n## location正则写法\r\n\r\n```sh\r\nlocation = / {\r\n    # 精确匹配 / ，主机名后面不能带任何字符串\r\n    [ configuration A ]\r\n}\r\n\r\nlocation / {\r\n    # 因为所有的地址都以 / 开头，所以这条规则将匹配到所有请求\r\n    # 但是正则和最长字符串会优先匹配\r\n    [ configuration B ]\r\n}\r\n\r\nlocation /documents/ {\r\n    # 匹配任何以 /documents/ 开头的地址，匹配符合以后，还要继续往下搜索\r\n    # 只有后面的正则表达式没有匹配到时，这一条才会采用这一条\r\n    [ configuration C ]\r\n}\r\n\r\nlocation ~ /documents/Abc {\r\n    # 匹配任何以 /documents/ 开头的地址，匹配符合以后，还要继续往下搜索\r\n    # 只有后面的正则表达式没有匹配到时，这一条才会采用这一条\r\n    [ configuration CC ]\r\n}\r\n\r\nlocation ^~ /images/ {\r\n    # 匹配任何以 /images/ 开头的地址，匹配符合以后，停止往下搜索正则，采用这一条。\r\n    [ configuration D ]\r\n}\r\n\r\nlocation ~* \\.(gif|jpg|jpeg)$ {\r\n    # 匹配所有以 gif,jpg或jpeg 结尾的请求\r\n    # 然而，所有请求 /images/ 下的图片会被 config D 处理，因为 ^~ 到达不了这一条正则\r\n    [ configuration E ]\r\n}\r\n\r\nlocation /images/ {\r\n    # 字符匹配到 /images/，继续往下，会发现 ^~ 存在\r\n    [ configuration F ]\r\n}\r\n\r\nlocation /images/abc {\r\n    # 最长字符匹配到 /images/abc，继续往下，会发现 ^~ 存在\r\n    # F与G的放置顺序是没有关系的\r\n    [ configuration G ]\r\n}\r\n\r\nlocation ~ /images/abc/ {\r\n    # 只有去掉 config D 才有效：先最长匹配 config G 开头的地址，继续往下搜索，匹配到这一条正则，采用\r\n    [ configuration H ]\r\n}\r\n\r\nlocation ~* /js/.*/\\.js\r\n```\r\n\r\n> - `=` 开头表示精确匹配（如 A 中只匹配根目录结尾的请求，后面不能带任何字符串）\r\n> - `^~` 开头表示uri以某个常规字符串开头，不是正则匹配\r\n> - `~` 开头表示区分大小写的正则匹配\r\n> - `~*` 开头表示不区分大小写的正则匹配\r\n> - `/` 通用匹配，如果没有其它匹配，任何请求都会匹配到\r\n\r\n> 顺序 no优先级：\r\n> (location =) > (location 完整路径) > (location ^~ 路径) > (location ~,~* 正则顺序) > (location 部分起始路径) > (/)\r\n\r\n上面的匹配结果\r\n按照上面的location写法，以下的匹配示例成立：\r\n\r\n- / -> config A\r\n精确完全匹配，即使/index.html也匹配不了\r\n- /downloads/download.html -> config B\r\n匹配B以后，往下没有任何匹配，采用B\r\n- /images/1.gif -> configuration D\r\n匹配到F，往下匹配到D，停止往下\r\n- /images/abc/def -> config D\r\n最长匹配到G，往下匹配D，停止往下\r\n你可以看到 任何以/images/开头的都会匹配到D并停止，FG写在这里是没有任何意义的，H是永远轮不到的，这里只是为了说明匹配顺序\r\n- /documents/document.html -> config C\r\n匹配到C，往下没有任何匹配，采用C\r\n- /documents/1.jpg -> configuration E\r\n匹配到C，往下正则匹配到E\r\n- /documents/Abc.jpg -> config CC\r\n最长匹配到C，往下正则顺序匹配到CC，不会往下到E\r\n\r\n## 实际使用建议\r\n\r\n实际使用中，个人觉得至少有三个匹配规则定义，如下：\r\n\r\n```sh\r\n#直接匹配网站根，通过域名访问网站首页比较频繁，使用这个会加速处理，官网如是说。\r\n#这里是直接转发给后端应用服务器了，也可以是一个静态首页\r\n# 第一个必选规则\r\nlocation = / {\r\n    proxy_pass http://tomcat:8080/index\r\n}\r\n# 第二个必选规则是处理静态文件请求，这是nginx作为http服务器的强项\r\n# 有两种配置模式，目录匹配或后缀匹配,任选其一或搭配使用\r\nlocation ^~ /static/ {\r\n    root /webroot/static/;\r\n}\r\nlocation ~* \\.(gif|jpg|jpeg|png|css|js|ico)$ {\r\n    root /webroot/res/;\r\n}\r\n#第三个规则就是通用规则，用来转发动态请求到后端应用服务器\r\n#非静态文件请求就默认是动态请求，自己根据实际把握\r\n#毕竟目前的一些框架的流行，带.php,.jsp后缀的情况很少了\r\nlocation / {\r\n    proxy_pass http://tomcat:8080/\r\n}\r\n```\r\n\r\n## Rewrite规则\r\n\r\nrewrite功能就是，使用nginx提供的全局变量或自己设置的变量，结合正则表达式和标志位实现url重写以及重定向。\r\nrewrite只能放在 server{}, location{}, if{} 中，并且只能对域名后边的除去传递的参数外的字符串起作用，例如： http://seanlook.com/a/we/index.php?id=1&u=str 只对 /a/we/index.php 重写。\r\n语法：rewrite regex replacement [flag];\r\n\r\n如果相对域名或参数字符串起作用，可以使用全局变量匹配，也可以使用proxy_pass反向代理。\r\n\r\n表明看rewrite和location功能有点像，都能实现跳转，主要区别在于rewrite是在同一域名内更改获取资源的路径，而location是对一类路径做控制访问或反向代理，可以proxy_pass到其他机器。很多情况下rewrite也会写在location里，它们的执行顺序是：\r\n\r\n1. 执行server块的rewrite指令\r\n2. 执行location匹配\r\n3. 执行选定的location中的rewrite指令\r\n\r\n如果其中某步URI被重写，则重新循环执行1-3，直到找到真实存在的文件；循环超过10次，则返回500 Internal Server Error错误。\r\n\r\n### 1、flag 标志位\r\n\r\n- last : 相当于Apache的[L]标记，表示完成rewrite\r\n- break : 停止执行当前虚拟主机的后续rewrite指令集\r\n- redirect : 返回302临时重定向，地址栏会显示跳转后的地址\r\n- permanent : 返回301永久重定向，地址栏会显示跳转后的地址\r\n\r\n因为301和302不能简单的只返回状态码，还必须有重定向的URL，这就是return指令无法返回301, 302的原因了。这里 last 和 break 区别有点难以理解：\r\n\r\n1. last一般写在server和if中，而break一般使用在location中\r\n2. last不终止重写后的url匹配，即新的url会再从server走一遍匹配流程，而break终止重写后的匹配\r\n3. break和last都能组织继续执行后面的rewrite指令\r\n\r\n### 2、if 指令\r\n\r\n语法为if(condition){...}，对给定的条件condition进行判断。如果为真，大括号内的rewrite指令将被执行，if条件(conditon)可以是如下任何内容：\r\n\r\n- 当表达式只是一个变量时，如果值为空或任何以0开头的字符串都会当做false\r\n- 直接比较变量和内容时，使用 `=` 或 `!=`\r\n- `~` 正则表达式匹配，`~*` 不区分大小写的匹配，`!~` 区分大小写的不匹配\r\n\r\n`-f` 和 `!-f` 用来判断是否存在文件\r\n`-d` 和 `!-d` 用来判断是否存在目录\r\n`-e` 和 `!-e` 用来判断是否存在文件或目录\r\n`-x` 和 `!-x` 用来判断文件是否可执行\r\n\r\n例如：\r\n\r\n```sh\r\nif ($http_user_agent ~ MSIE) {\r\n    rewrite ^(.*)$ /msie/$1 break;\r\n} #如果UA包含\"MSIE\"，rewrite请求到/msid/目录下\r\n\r\nif ($http_cookie ~* \"id=([^;]+)(?:;|$)\") {\r\n    set $id $1;\r\n} #如果cookie匹配正则，设置变量$id等于正则引用部分\r\n\r\nif ($request_method = POST) {\r\n    return 405;\r\n} #如果提交方法为POST，则返回状态405（Method not allowed）。return不能返回301,302\r\n\r\nif ($slow) {\r\n    limit_rate 10k;\r\n} #限速，$slow可以通过 set 指令设置\r\n\r\nif (!-f $request_filename){\r\n    break;\r\n    proxy_pass  http://127.0.0.1;\r\n} #如果请求的文件名不存在，则反向代理到localhost 。这里的break也是停止rewrite检查\r\n\r\nif ($args ~ post=140){\r\n    rewrite ^ http://example.com/ permanent;\r\n} #如果query string中包含\"post=140\"，永久重定向到example.com\r\n\r\nlocation ~* \\.(gif|jpg|png|swf|flv)$ {\r\n    valid_referers none blocked www.jefflei.com www.leizhenfang.com;\r\n    if ($invalid_referer) {\r\n        return 404;\r\n    } #防盗链\r\n}\r\n```\r\n\r\n### 3、全局变量\r\n\r\n下面是可以用作if判断的全局变量\r\n\r\n| 变量名称 | 说明 |\r\n| --- | --- |\r\n| $args | 这个变量等于请求行中的参数，同$query_string |\r\n| $content_length | 请求头中的Content-length字段 |\r\n| $content_type | 请求头中的Content-Type字段 |\r\n| $document_root | 当前请求在root指令中指定的值 |\r\n| $host | 请求主机头字段，否则为服务器名称 |\r\n| $http\\_user_agent | 客户端agent信息 |\r\n| $http_cookie | 客户端cookie信息 |\r\n| $limit_rate | 这个变量可以限制连接速率 |\r\n| $request_method | 客户端请求的动作，通常为GET或POST |\r\n| $remote_addr | 客户端的IP地址 |\r\n| $remote_port | 客户端的端口 |\r\n| $remote_user | 已经经过Auth Basic Module验证的用户名 |\r\n| $request_filename | 当前请求的文件路径，由root或alias指令与URI请求生成 |\r\n| $scheme | HTTP方法（如http，https） |\r\n| $server_protocol | 请求使用的协议，通常是HTTP/1.0或HTTP/1.1 |\r\n| $server_addr | 服务器地址，在完成一次系统调用后可以确定这个值 |\r\n| $server_name | 服务器名称 |\r\n| $server_port | 请求到达服务器的端口号 |\r\n| $request_uri | 包含请求参数的原始URI，不包含主机名，如：\"/foo/bar.php?arg=baz\" |\r\n| $uri | 不带请求参数的当前URI，$uri不包含主机名，如：\"/foo/bar.html\" |\r\n| $document_uri | 与$uri相同 |\r\n\r\n```sh\r\n# 例：http://localhost:88/test1/test2/test.php\r\n$host：localhost\r\n$server_port：88\r\n$request_uri：http://localhost:88/test1/test2/test.php\r\n$document_uri：/test1/test2/test.php\r\n$document_root：/var/www/html\r\n$request_filename：/var/www/html/test1/test2/test.php\r\n```\r\n\r\n### 4、常用正则\r\n\r\n| 匹配符 | 说明 |\r\n| --- | --- |\r\n| `.` | 匹配除换行符以外的任意字符 |\r\n| `?` | 重复0次或1次 |\r\n| `+` | 重复1次或更多次 |\r\n| `*` | 重复0次或更多次 |\r\n| `\\d` | 匹配数字 |\r\n| `^` | 匹配字符串的开始 |\r\n| `$` | 匹配字符串的介绍 |\r\n| `{n}` | 重复n次 |\r\n| `{n,}` | 重复n次或更多次 |\r\n| `[c]` | 匹配单个字符c |\r\n| `[a-z]` | 匹配a-z小写字母的任意一个 |\r\n\r\n> 小括号()之间匹配的内容，可以在后面通过$1来引用，$2表示的是前面第二个()里的内容。\r\n> 正则里面容易让人困惑的是 `\\` 转义特殊字符。\r\n\r\n### 5、rewrite实例\r\n\r\n```sh\r\n# 例1：\r\nhttp {\r\n    # 定义image日志格式\r\n    log_format imagelog \'[$time_local] \' $image_file \' \' $image_type \' \' $body_bytes_sent \' \' $status;\r\n    # 开启重写日志\r\n    rewrite_log on;\r\n\r\n    server {\r\n        root /home/www;\r\n\r\n        location / {\r\n            # 重写规则信息\r\n            error_log logs/rewrite.log notice;\r\n            # 注意这里要用‘’单引号引起来，避免{}\r\n            rewrite \'^/images/([a-z]{2})/([a-z0-9]{5})/(.*)\\.(png|jpg|gif)$\' /data?file=$3.$4;\r\n            # 注意不能在上面这条规则后面加上“last”参数，否则下面的set指令不会执行\r\n            set $image_file $3;\r\n            set $image_type $4;\r\n        }\r\n\r\n        location /data {\r\n            # 指定针对图片的日志格式，来分析图片类型和大小\r\n            access_log logs/images.log mian;\r\n            root /data/images;\r\n            # 应用前面定义的变量。判断首先文件在不在，不在再判断目录在不在，如果还不在就跳转到最后一个url里\r\n            try_files /$arg_file /image404.html;\r\n        }\r\n        location = /image404.html {\r\n            # 图片不存在返回特定的信息\r\n            return 404 \"image not found\\n\";\r\n        }\r\n}\r\n```\r\n\r\n> 对形如 /images/ef/uh7b3/test.png 的请求，重写到 /data?file=test.png，于是匹配到 location /data，先看 /data/images/test.png 文件存不存在，如果存在则正常响应，如果不存在则重写tryfiles到新的image404 location，直接返回404状态码。\r\n\r\n```sh\r\n# 例2：\r\nrewrite ^/images/(.*)_(\\d+)x(\\d+)\\.(png|jpg|gif)$ /resizer/$1.$4?width=$2&height=$3? last;\r\n```\r\n\r\n> 对形如/images/bla_500x400.jpg的文件请求，重写到 /resizer/bla.jpg?width=500&height=400 地址，并会继续尝试匹配location。', '2017-07-30 18:37:59', '2017-07-30 18:37:59');
INSERT INTO `go_doc` VALUES (2, 'golang test测试使用', 1, 2, 'golang,test,测试', '# golang test测试使用\r\n\r\n## 创建需要测试的文件mysql.go\r\n\r\n```go\r\npackage mysql\r\n\r\nimport (\r\n    \"database/sql\"\r\n    _ \"github.com/go-sql-driver/mysql\"\r\n)\r\n\r\nfunc findByPk(pk int) (int, error) {\r\n    var num int = 0\r\n\r\n    db, err := sql.Open(\"mysql\", \"root:@tcp(127.0.0.1:3306)/plugin_master?charset=utf8\")\r\n\r\n    if err != nil {\r\n        return 0, err\r\n    }\r\n\r\n    defer db.Close()\r\n\r\n    stmtOut, err := db.Prepare(\"select id from t_admin where id=?\")\r\n\r\n    if err != nil {\r\n        return 0, err\r\n    }\r\n\r\n    defer stmtOut.Close()\r\n\r\n    err = stmtOut.QueryRow(pk).Scan(&num)\r\n\r\n    if err != nil {\r\n        return 0, err\r\n    }\r\n\r\n    return num, nil\r\n}\r\n```\r\n\r\n## 创建单元测试用例文件mysql_test.go\r\n\r\n文件名必须是 `*_test.go` 的类型，`*` 代表要测试的文件名，函数名必须以 `Test` 开头如：`TestXxx` 或 `Test_xxx`\r\n\r\n```go\r\npackage mysql\r\n\r\nimport (\r\n    \"testing\"\r\n)\r\n\r\nfunc Test_findByPk(t *testing.T) {\r\n    num, err := findByPk(1)\r\n\r\n    if err != nil {\r\n        t.Errorf(\"mysql find error: %s\", err.Error())\r\n    }\r\n\r\n    t.Log(num)\r\n}\r\n```\r\n\r\n测试所有的文件 `go test`，将对当前目录下的所有 `*_test.go` 文件进行编译并自动运行测试\r\n\r\n测试某个文件使用 `-file` 参数：`go test –file *.go`\r\n例如：`go test -file mysql_test.go`\r\n\r\n> 注：`-file` 参数不是必须的，可以省略，如果你输入 `go test b_test.go` 也会得到一样的效果\r\n\r\n测试某个方法 `go test -run=\"Test_xxx\"`\r\n\r\n`-v` 参数：`go test -v ...` 表示无论用例是否测试通过都会显示结果，不加 `v` 表示只显示未通过的用例结果\r\n\r\n## 创建benchmark性能测试用例文件mysql_b_test.go\r\n\r\n文件名必须是 `*_b_test.go` 的类型，`*` 代表要测试的文件名，函数名必须以 `Benchmark` 开头如：`BenchmarkXxx` 或`Benchmark_xxx`\r\n\r\n```go\r\npackage mysql\r\n\r\nimport (\r\n    \"testing\"\r\n)\r\n\r\nfunc Benchmark_findByPk(b *testing.B) {\r\n    for i := 0; i < b.N; i++ { //use b.N for looping\r\n        findByPk(1)\r\n    }\r\n}\r\n```\r\n\r\n进行所有go文件的benchmark测试：`go test -bench=\".*\"` 或 `go test . -bench=\".*\"`\r\n\r\n对某个go文件进行benchmark测试：`go test mysql_b_test.go -bench=\".*\"`\r\n\r\n## 用性能测试生成CPU状态图(暂未测试使用)\r\n\r\n使用命令：\r\n\r\n```sh\r\ngo test -bench=\".*\" -cpuprofile=cpu.prof -c\r\n```\r\n\r\n> cpuprofile是表示生成的cpu profile文件\r\n> -c是生成可执行的二进制文件，这个是生成状态图必须的，它会在本目录下生成可执行文件mysql.test\r\n\r\n然后使用go tool pprof工具：\r\n\r\n```sh\r\ngo tool pprof mysql.test cpu.prof\r\n```\r\n\r\n> 调用web（需要安装graphviz）来生成svg文件，生成后使用浏览器查看svg文件\r\n> 参考 http://www.cnblogs.com/yjf512/archive/2013/01/18/2865915.html', '2017-07-30 18:41:09', '2017-07-30 18:41:09');
INSERT INTO `go_doc` VALUES (3, 'golang nginx 配置', 1, 1, 'golang,nginx', '# golang nginx 配置\r\n\r\n```sh\r\nserver{\r\n    listen 80;\r\n    server_name api.golang.com;\r\n    access_log /data/logs/nginx/access/api.golang.com.log;\r\n    error_log /data/logs/nginx/error/api.golang.com.log;\r\n\r\n    location / {\r\n        proxy_set_header X-Forwarded-For $remote_addr;\r\n        proxy_set_header Host $http_host;\r\n        proxy_http_version 1.1;\r\n        proxy_set_header Connection \"\";\r\n\r\n        proxy_pass http://http_backend;\r\n    }\r\n}\r\n\r\nupstream http_backend {\r\n    server 127.0.0.1:8000;\r\n    keepalive 256;\r\n}\r\n```', '2017-07-30 18:42:04', '2017-07-30 18:42:04');
INSERT INTO `go_doc` VALUES (4, '用户详情', 3, 3, 'users,用户', '# 用户详情\r\n\r\n### 请求和参数说明\r\n\r\n| 请求 | /users  | GET  | v1  |\r\n| --- | --- | --- |\r\n| 参数  | 是否为空 | 说明  | 缺省值  |\r\n| id  | 否 | 用户ID  | 无 |\r\n\r\n### 返回\r\n\r\n```json\r\n{\r\n	\"id\": 1,\r\n	\"name\": \"demo\",\r\n	\"email\": \"demo@demo.com\",\r\n	\"role\": 1，\r\n	\"status\": 1\r\n}\r\n```', '2017-07-30 18:51:59', '2017-07-30 18:51:59');
COMMIT;

-- ----------------------------
-- Table structure for go_history
-- ----------------------------
DROP TABLE IF EXISTS `go_history`;
CREATE TABLE `go_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `category_id` int(11) NOT NULL COMMENT '类别ID',
  `project_id` int(11) NOT NULL COMMENT '项目ID',
  `doc_id` int(11) NOT NULL COMMENT '文档ID',
  `flag` tinyint(4) NOT NULL COMMENT '类别(1：创建；2：修改)',
  `created_at` datetime NOT NULL COMMENT '添加时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_doc` (`doc_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COMMENT='操作历史表';

-- ----------------------------
-- Records of go_history
-- ----------------------------
BEGIN;
INSERT INTO `go_history` VALUES (1, 1, 1, 1, 1, 1, '2017-07-30 18:37:59', '2017-07-30 18:37:59');
INSERT INTO `go_history` VALUES (2, 1, 1, 2, 2, 1, '2017-07-30 18:41:09', '2017-07-30 18:41:09');
INSERT INTO `go_history` VALUES (3, 1, 1, 1, 3, 1, '2017-07-30 18:42:04', '2017-07-30 18:42:04');
INSERT INTO `go_history` VALUES (4, 1, 3, 3, 4, 1, '2017-07-30 18:51:59', '2017-07-30 18:51:59');
COMMIT;

-- ----------------------------
-- Table structure for go_project
-- ----------------------------
DROP TABLE IF EXISTS `go_project`;
CREATE TABLE `go_project` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) NOT NULL COMMENT '名称',
  `category_id` int(11) NOT NULL COMMENT '类别ID',
  `description` varchar(255) NOT NULL COMMENT '描述',
  `created_at` datetime NOT NULL COMMENT '添加时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `index_category` (`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='项目表';

-- ----------------------------
-- Records of go_project
-- ----------------------------
BEGIN;
INSERT INTO `go_project` VALUES (1, 'nginx', 1, '记录有关nginx的相关技术知识', '2017-07-30 18:36:55', '2017-07-30 18:36:55');
INSERT INTO `go_project` VALUES (2, 'go', 1, '记录golang的相关开发知识', '2017-07-30 18:40:05', '2017-07-30 18:40:05');
INSERT INTO `go_project` VALUES (3, '测试项目', 3, '这是一个测试项目', '2017-07-30 18:45:25', '2017-07-30 18:45:25');
COMMIT;

-- ----------------------------
-- Table structure for go_user
-- ----------------------------
DROP TABLE IF EXISTS `go_user`;
CREATE TABLE `go_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) NOT NULL COMMENT '用户名',
  `email` varchar(50) NOT NULL COMMENT '邮箱',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `salt` varchar(20) NOT NULL COMMENT '加密盐',
  `role` int(11) NOT NULL COMMENT '角色',
  `last_login_ip` varchar(20) DEFAULT '' COMMENT '最近登录IP',
  `last_login_time` datetime DEFAULT '1970-01-01 00:00:00' COMMENT '最近登录时间',
  `created_at` datetime NOT NULL COMMENT '添加时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_email` (`email`),
  UNIQUE KEY `index_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- ----------------------------
-- Records of go_user
-- ----------------------------
BEGIN;
INSERT INTO `go_user` VALUES (1, 'admin', 'admin@qq.com', '0de3506766e559a5b0344284709fd0ad', '2ngW5iCS2XpAezy9', 3, '::1', '2017-07-30 16:28:41', '2017-06-04 21:03:19', '2017-07-30 16:28:41');
INSERT INTO `go_user` VALUES (2, 'demo', 'demo@qq.com', '027a94619ce748fac471a905af271894', 'QAfY0TJDhHHmm%8R', 2, '::1', '2017-07-30 16:20:20', '2017-06-13 15:14:49', '2017-07-30 18:29:32');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/vladpereskokov/Technopark_HighLoad-nginx/src/constants"
	"github.com/vladpereskokov/Technopark_HighLoad-nginx/src/models"
	"net"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Connection net.Conn
	Request    *models.Request
	Response   *models.Response
	Dir        string
}

func (handler *Handler) Start(channel chan net.Conn) {
	for {
		conn := <-channel
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		conn.Write([]byte("Message received."))
		conn.Close()
	}
}

func (handler *Handler) get_path() string {
	return handler.Request.GetPath()
}

func (handler *Handler) set_path(new_path string) {
	handler.Request.SetPath(new_path)
}

func (handler *Handler) set_header(key string, value string) {
	handler.Response.Headers[key] = value
}

func (handler *Handler) set_status(status int) {
	handler.Response.SetStatus(status)
}

func (handler *Handler) read_request() {
	buffer := make([]byte, 1024)
	_, err := handler.Connection.Read(buffer)
	if err != nil {
		fmt.Println("Read request error ", err)
	}
	raw_request := string(buffer[:bytes.Index(buffer, []byte{0})])
	start_string := strings.Split(raw_request, constants.STRING_SEPARATOR)[0]
	fmt.Println(start_string)
	handler.parse_start_string(start_string)
}

func (handler *Handler) parse_start_string(start_string string) {
	splited_string := strings.Split(start_string, " ")
	if len(splited_string) != 3 {
		handler.set_status(400)
		return
	}
	handler.Request.Method.SetMethod(splited_string[0])
	parsed_url, err := url.Parse(splited_string[1])
	if err != nil || !strings.HasPrefix(splited_string[2], "HTTP/") {
		handler.set_status(400)
	}
	handler.Request.Url = parsed_url
}

func (handler *Handler) process_request() {
	if !handler.Response.IsOk() {
		handler.set_content_headers(nil)
		return
	}
	if !contains(constants.IMPLEMENTED_METHODS, handler.Request.Method.GetMethod()) {
		handler.set_status(405)
	} else {
		handler.preprocess_path()
	}
}

//preproccess path and check file errors
func (handler *Handler) preprocess_path() {
	//handler.set_path(handler.Factory.root + handler.get_path())
	file_info := handler.check_path(false)
	if file_info != nil && file_info.IsDir() {
		handler.set_path(handler.get_path() + constants.INDEX_FILE)
		file_info = handler.check_path(true)
	}
	handler.set_content_headers(file_info)
}

func (handler *Handler) check_path(is_dir bool) os.FileInfo {
	request_path := handler.get_path()
	clear_path := path.Clean(request_path)
	handler.set_path(clear_path)
	info, err := os.Stat(request_path)
	if err != nil {

		if os.IsNotExist(err) && !is_dir {
			handler.set_status(404)
		} else {
			handler.set_status(403)
		}
	}
	//} else if !strings.Contains(clear_path, handler.Factory.root) {
	//	handler.set_status("forbidden")
	//}
	return info
}

func (handler *Handler) set_content_headers(info os.FileInfo) {
	if handler.Response.IsOk() {
		handler.set_header("Content-Length", strconv.Itoa(int(info.Size())))
		handler.set_header("Content-Type", handler.get_content_type())
	} else {
		handler.set_header("Content-Length", strconv.Itoa(len(handler.get_error_body())))
		handler.set_header("Content-Type", constants.ERROR_BODY_MIME_TYPE)
	}
}

func contains(arr []string, value string) bool {
	for _, elem := range arr {
		if elem == value {
			return true
		}
	}
	return false
}

func (handler *Handler) get_content_type() string {
	extension := ""
	request_path := handler.get_path()
	last_dot := strings.LastIndex(request_path, ".")
	if last_dot >= 0 {
		extension = request_path[last_dot:]
	}
	val, ok := constants.CONTENT_TYPES[extension]
	if ok {
		return val
	} else {
		return constants.DEFAULT_MIME_TYPE
	}
}

func (handler *Handler) clear() {
	//handler.Factory = nil
	handler.Connection.Close()
}

func (handler Handler) write_response() {
	handler.write_string(constants.HTTP_VERSION + " " + handler.Response.Status.Message)
	handler.write_headers()
	handler.write_string("") // empty string after headers
	if handler.Request.Method.GetMethod() != "HEAD" {
		handler.write_body()
	}
	fmt.Println(handler.Request.Method, " ", handler.get_path(), " ", handler.Response.Status.Code)
}

func (handler *Handler) write_string(str string) {
	handler.Connection.Write([]byte(str + constants.STRING_SEPARATOR))
}

func (handler *Handler) write_body() {
	if handler.Response.IsOk() {
		handler.write_ok_body()
	} else {
		handler.write_error_body()
	}
}

func (handler *Handler) write_ok_body() {
	file, err := os.Open(handler.get_path())
	if err != nil {
		fmt.Println("Can't open file ", handler.get_path())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_, read_err := reader.WriteTo(handler.Connection)
	if read_err != nil {
		fmt.Println("Some error on read or write file ", handler.get_path())
	}
}

func (handler *Handler) write_error_body() {
	body := []byte(handler.get_error_body())
	handler.Connection.Write(body)
}

func (handler *Handler) get_error_body() string {
	body := "<html><body><h1>"
	body += handler.Response.Status.Message
	body += "</h1></body></html>"
	return body
}

func (handler Handler) write_headers() {
	handler.write_common_headers()
	handler.write_specific_headers()
}

func (handler Handler) write_common_headers() {
	handler.write_string("Date: " + time.Now().String())
	handler.write_string("Server: " + constants.SERVER)
	handler.write_string("Connection: close")
}

func (handler Handler) write_specific_headers() {
	for key, value := range handler.Response.Headers {
		handler.write_string(key + ": " + value)
	}
}

func (handler *Handler) Handle() {
	handler.read_request()
	handler.process_request()
	handler.write_response()
	handler.clear()
}

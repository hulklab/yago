package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func ExecuteTemplate(tmplText string, data interface{}) string {
	if len(tmplText) == 0 {
		panic("tmplText is required")
	}

	t, err := template.New("new").Parse(tmplText)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer

	err = t.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	return b.String()
}

type RpcTmplData struct {
	Package string
	Name    string
}

var RpcTemplate = `package {{.Package}} 

import (
	"context"
	"log"

	"github.com/hulklab/yago"

	pb "github.com/hulklab/yago/example/app/modules/home/homerpc/homepb"
)

type {{.Name}}Rpc struct {
}

func init() {
	h := new({{.Name}}Rpc)
	pb.RegisterHomeServer(yago.RpcServer, h)
}

func (r *{{.Name}}Rpc) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Data: "Hello " + in.Name}, nil
}
`

type CmdTmplData struct {
	Package string
	Name    string
}

var CmdTemplate = `package {{.Package}} 

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basecmd"
	"github.com/spf13/cobra"
)

type {{.Name}}Cmd struct {
	basecmd.BaseCmd
}

func init() {
	c := new({{.Name}}Cmd)
	// 注册路由
	yago.AddCmdRouter("demo", "Demo action", c.DemoAction, yago.CmdStringArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: true,
	})
}

func (c *{{.Name}}Cmd) DemoAction(cmd *cobra.Command, args []string) {

}

`

type TaskTmplData struct {
	Package string
	Name    string
}

var TaskTemplate = `package {{.Package}} 

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basetask"
)

type {{.Name}}Task struct {
	basetask.BaseTask
}

func init() {
	t := new({{.Name}}Task)
	yago.AddTaskRouter("@loop", t.HelloAction)
	yago.AddTaskRouter("0 */1 * * * *", t.HelloAction)
}

func (t *{{.Name}}Task) HelloAction() {
	//t.RunLoop(func() {
	//})
}
`

type ApiTmplData struct {
	Package string
	LName   string
	OName   string
}

var ApiTemplate = `package {{.Package}}

import (
	"fmt"
	"log"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basethird"
	"github.com/levigross/grequests"
)

type {{.LName}}Api struct {
	basethird.HttpThird
}

func Ins() *{{.LName}}Api{
	name := "{{.OName}}_api"
	v := yago.Component.Ins(name, func() interface{} {
		api := new({{.LName}}Api)

		err := api.InitConfig(name)
		if err != nil {
			log.Fatal("init {{.OName}} api config error:", err.Error())
		}
		return api
	})
	return v.(*{{.LName}}Api)
}

func (a *{{.LName}}Api) Hello() {

	ro := &grequests.RequestOptions{
		JSON: map[string]interface{}{},
	}

	resp, err := a.Post("/hello", nil, ro)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.String())
	}

}
`

// ---------------------- http-method -----------------------

type HttpMethodTmplData struct {
	HasReq     bool
	HasResp    bool
	HasErr     bool
	CamelName  string
	Method     string
	ModuleName string
	ReqPackage string
	ReqName    string
}

func (d HttpMethodTmplData) RecvStr() string {
	if d.HasResp && d.HasErr {
		return "resp,err"
	} else if !d.HasResp && d.HasErr {
		return "err"
	} else if d.HasResp && !d.HasErr {
		return "resp"
	}

	return ""
}

func (d HttpMethodTmplData) ReturnStr() string {
	if d.HasResp && d.HasErr {
		return "c.SetDataOrErr(resp,err)"
	} else if !d.HasResp && d.HasErr {
		return "c.SetDataOrErr(g.Hash{},err)"
	} else if d.HasResp && !d.HasErr {
		return "c.SetDataOrErr(resp,nil)"
	}

	return ""
}

var HttpMethodTemplate = `
func (h *{{.CamelName}}Http) {{.Method}}Action(c *yago.Ctx) {
	ctx := h.GetTraceCtx(c)

	{{if .HasReq}}
	req := &{{.ReqPackage}}.{{.ReqName}}{}

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	{{.RecvStr}} := {{.ModuleName}}service.New{{.CamelName}}Service(ctx).{{.Method}}(req)
	{{else}}
	{{.RecvStr}} := {{.ModuleName}}service.New{{.CamelName}}Service(ctx).{{.Method}}()
	{{end}}

	{{.ReturnStr}}
}
`

type HttpRouteTmplData struct {
	Group      string
	ModuleName string
	LispName   string
	LispMethod string
	Method     string
}

func (d HttpRouteTmplData) GroupName() string {
	if len(d.Group) == 0 {
		d.Group = "Root"
	}

	return d.Group
}

var HttpRouteTemplate = `	ghttp.{{.GroupName}}.Post("/{{.ModuleName}}/{{.LispName}}/{{.LispMethod}}", h.{{.Method}}Action)`

// ------------------------- http ----------------------------
type HttpTmplData struct {
	Empty          bool
	Package        string
	Name           string
	DtoPackage     string
	ServicePackage string
	ModName        string
	ModuleName     string
	AddRoute       string
	DelRoute       string
	UpdateRoute    string
	ListRoute      string
	DetailRoute    string
	Entry          string
}

func (d HttpTmplData) BaseHttpImport() string {
	return fmt.Sprintf("%s/app/g/ghttp", d.ModName)
}
func (d HttpTmplData) BaseHttpName() string {
	return "ghttp.BaseHttp"
}

var HttpTemplate = `package {{.Package}}
import (
"{{.BaseHttpImport}}"
"github.com/hulklab/yago"

{{if not .Empty}}
"{{.ModName}}/app/g"
"{{.ModName}}/app/modules/{{.ModuleName}}/{{.DtoPackage}}"
"{{.ModName}}/app/modules/{{.ModuleName}}/{{.ServicePackage}}"
{{end}}
)

type {{.Name}}Http struct {
	{{.BaseHttpName}}
}

func init() {
	h := new({{.Name}}Http)

	{{if .Empty}}
	_ = h
	{{else}}
	{{.AddRoute}}
	{{.DelRoute}}
	{{.UpdateRoute}}
	{{.ListRoute}}
	{{.DetailRoute}}
	{{end}}
}

{{if not .Empty}}
func (h *{{.Name}}Http) ListAction(c *yago.Ctx) {
	req := &{{.DtoPackage}}.{{.Name}}ListReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req);err != nil {
		c.SetError(err)
		return
	}

	data,err := {{.ServicePackage}}.New{{.Name}}Service(ctx).GetList(req)
	c.SetDataOrErr(data,err)
}

func (h *{{.Name}}Http) AddAction(c *yago.Ctx) {
	req := &{{.DtoPackage}}.{{.Name}}AddReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req);err != nil {
		c.SetError(err)
		return
	}

	data,err := {{.ServicePackage}}.New{{.Name}}Service(ctx).Add{{.Name}}(req)
	c.SetDataOrErr(data,err)
}

func (h *{{.Name}}Http) UpdateAction(c *yago.Ctx) {
	req := &{{.DtoPackage}}.{{.Name}}UpdateReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req);err != nil {
		c.SetError(err)
		return
	}

	err := {{.ServicePackage}}.New{{.Name}}Service(ctx).UpdateById(req)
	c.SetDataOrErr(g.Hash{},err)
}

func (h *{{.Name}}Http) DeleteAction(c *yago.Ctx) {
	req := &{{.DtoPackage}}.{{.Name}}DeleteReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req);err != nil {
		c.SetError(err)
		return
	}

	err := {{.ServicePackage}}.New{{.Name}}Service(ctx).DeleteById(req)
	c.SetDataOrErr(g.Hash{},err)
}

func (h *{{.Name}}Http) DetailAction(c *yago.Ctx) {
	req := &{{.DtoPackage}}.{{.Name}}DetailReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req);err != nil {
		c.SetError(err)
		return
	}

	data,err := {{.ServicePackage}}.New{{.Name}}Service(ctx).GetDetail(req)
	c.SetDataOrErr(data,err)
}
{{end}}
`

// ------------------------- model ------------------------
type ModelTmplData struct {
	Empty         bool
	Package       string
	ModName       string
	DaoImportPath string
	Lname         string
	Name          string
	DaoName       string
	DaoPackage    string
}

func (t ModelTmplData) GmodelPath() string {
	return fmt.Sprintf("%s/app/g/gmodel", t.ModName)
}

var ModelTemplate = `package {{.Package}}

import (
	"{{.GmodelPath}}"
	{{if not .Empty}}
	"{{.ModName}}/app/g"
	"{{.DaoImportPath}}"
	"fmt"
	{{end}}
)

type {{.Lname}}Model struct {
	gmodel.BaseModel

}

func New{{.Name}}Model(opts ...gmodel.Option) *{{.Lname}}Model {
	m := &{{.Lname}}Model{}
	m.Init(opts...)
	return m
}

{{if not .Empty}}
func (m *{{.Lname}}Model) tableName()  string{
	return new({{.DaoPackage}}.{{.DaoName}}).TableName()
}

func (m *{{.Lname}}Model) InsertOne(dao *{{.DaoPackage}}.{{.DaoName}})  (err error){
	_,err = m.GetSession().InsertOne(dao)

	return
}

func (m *{{.Lname}}Model) GetById(id int64)  (b bool,dao *{{.DaoPackage}}.{{.DaoName}},err error){
	dao = &{{.DaoPackage}}.{{.DaoName}}{}

	b,err = m.GetSession().Where("id = ?", id).Get(dao)
	if err != nil  || !b{
		return b,nil,err
	}

	return b,dao,err
}

func (m *{{.Lname}}Model) MustGetById(id int64)  (dao *{{.DaoPackage}}.{{.DaoName}},err error){
	dao = &{{.DaoPackage}}.{{.DaoName}}{}

	b,err := m.GetSession().Where("id = ?", id).Get(dao)
	if err != nil {
		return nil,err
	}

	if !b {
		return nil,fmt.Errorf("{{.Lname}} %d is not exist",id)
	}

	return dao,err
}

func (m *{{.Lname}}Model) GetListByIds(ids []int64) (list []*{{.DaoPackage}}.{{.DaoName}}, err error) {
	list = make([]*{{.DaoPackage}}.{{.DaoName}}, 0)

	err = m.GetSession().In("id", ids).Find(&list)

	return list, err
}

func (m *{{.Lname}}Model) GetAll()  (list []*{{.DaoPackage}}.{{.DaoName}},err error){
	list = make([]*{{.DaoPackage}}.{{.DaoName}},0)
	err = m.GetSession().Find(&list)

	return list,err
}

func (m *{{.Lname}}Model) GetAllMap()  (allMap map[int64]*{{.DaoPackage}}.{{.DaoName}},err error){
	allMap = make(map[int64]*{{.DaoPackage}}.{{.DaoName}})
	err = m.GetSession().Find(&allMap)

	return allMap,err
}

func (m *{{.Lname}}Model) DeleteById(id int64)  (err error){
	session := m.GetSession()
	session.Where("id=?",id)

	_,err = session.Delete(new({{.DaoPackage}}.{{.DaoName}}))
	return
}

func (m *{{.Lname}}Model) DeleteByIds(ids []int64)  (err error){
	session := m.GetSession()
	session.In("id",ids)

	_,err = session.Delete(new({{.DaoPackage}}.{{.DaoName}}))
	return
}


func (m *{{.Lname}}Model) UpdateAttrsById(id int64, attrs g.Hash) (err error) {
	if attrs == nil {
		return fmt.Errorf("attrs cannot be nil, id:%d", id)
	}

	_, err = m.GetSession().Table(m.tableName()).Where("id=?", id).Update(attrs)

	return
}
{{end}}
`

// ------------------------- service-method -------------------------
type MethodTmplData struct {
	StructName string
	ModuleName string
	Method     string
	HasReq     bool
	HasResp    bool
	HasErr     bool
}

func (d MethodTmplData) Req() string {
	if d.HasReq {
		return fmt.Sprintf("req *%sdto.%sReq", d.ModuleName, d.Method)
	}

	return ""
}

func (d MethodTmplData) Err() string {
	if d.HasErr {
		return "err error"
	}

	return ""
}

func (d MethodTmplData) Resp() string {
	if d.HasResp {
		return fmt.Sprintf("resp *%sdto.%sResp%s", d.ModuleName, d.Method, If(d.HasErr, ", "))
	}

	return ""
}

var MethodTemplate = `func (s *{{.StructName}}) {{.Method}}({{.Req}}) ({{.Resp}}{{.Err}}) {

	return
}
`

var DtoReqTemplate = `
type {{.Method}}Req struct {
	
}
`
var DtoRespTemplate = `
type {{.Method}}Resp struct {
	
}
`

// ------------------------- service ------------------
type ServiceTmplData struct {
	Empty            bool
	Package          string
	Lname            string
	Name             string
	ModName          string
	DtoPackage       string
	ModelPackage     string
	ModuleImportPath string
	DaoPackage       string
	DaoName          string
	DaoImportPath    string
	FieldMap         map[string]field
}

func (d ServiceTmplData) GetDetailAssign() string {
	structAssign := ""
	for _, field := range d.FieldMap {
		if field.NotShowJson() {
			continue
		}
		structAssign += fmt.Sprintf("\t\t%s:\tdao.%s,\n", field.Name, field.Name)
	}
	return structAssign
}

func (d ServiceTmplData) GetUpdateAssign() string {
	structAssign := ""
	for _, field := range d.FieldMap {
		if strings.ToLower(field.Name) == "id" {
			continue
		}
		if field.NotShowJson() || field.NotShowXorm() || field.IsOrmHook() {
			continue
		}
		structAssign += fmt.Sprintf("\t\t%s:\treq.%s,\n", field.Name, field.Name)
	}

	return structAssign
}

func (d ServiceTmplData) GetAddAssign() string {
	// 结构体赋值
	structAssign := ""
	for _, field := range d.FieldMap {
		if strings.ToLower(field.Name) == "id" {
			continue
		}

		if field.NotShowJson() || field.NotShowXorm() || field.IsOrmHook() {
			continue
		}
		structAssign += fmt.Sprintf("\t\t%s:\treq.%s,\n", field.Name, field.Name)
	}

	return structAssign
}

var ServiceTemplate = `package {{.Package}}

import (
	"{{.ModName}}/app/g/gservice"
	"{{.ModName}}/app/libs/trace"
	{{if not .Empty}}
	"{{.DaoImportPath}}"
	"{{.ModuleImportPath}}/{{.ModelPackage}}"
	"{{.ModuleImportPath}}/{{.DtoPackage}}"
	"{{.ModName}}/app/g/gmodel"
	"fmt"
	{{end}}
)

type {{.Lname}}Service struct {
	gservice.BaseService
}

func New{{.Name}}Service(ctx *trace.Context)  *{{.Lname}}Service{
	s := &{{.Lname}}Service{}
	s.Init(ctx)
	return s
}

{{if not .Empty}}
func (s *{{.Lname}}Service) GetList(req *{{.DtoPackage}}.{{.Name}}ListReq)  (resp *{{.DtoPackage}}.{{.Name}}ListResp,err error){
	resp = &{{.DtoPackage}}.{{.Name}}ListResp{
		List: make([]*{{.DtoPackage}}.{{.Name}},0),
	}

	session := {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).GetSession()

	list := make([]*{{.DaoPackage}}.{{.DaoName}},0)

	if len(req.Q) > 0 {
		// session.Where("? like ?","%"+req.Q+"%")
	}

	session.OrderBy("id desc")

	if req.PageNum > 0 && req.PageSize > 0 {
		session.Limit(req.PageSize, req.PageSize*(req.PageNum-1))
	}

	total,err := session.FindAndCount(&list)
	if err != nil {
		return
	}

	resp.Total = total

	for _, dao := range list {
		info := &{{.DtoPackage}}.{{.Name}}{
{{.GetDetailAssign}}
		}

		resp.List = append(resp.List,info)
	}
	return
}


func (s *{{.Lname}}Service) GetDetail(req *{{.DtoPackage}}.{{.Name}}DetailReq) (resp *{{.DtoPackage}}.{{.Name}}, err error) {

	dao, err := {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return nil, err
	}

	info := &{{.DtoPackage}}.{{.Name}}{
{{.GetDetailAssign}}
	}

	return info, nil
}

func (s *{{.Lname}}Service) Add{{.Name}}(req *{{.DtoPackage}}.{{.Name}}AddReq) (resp *{{.DtoPackage}}.{{.Name}}AddResp, err error) {

	dao := &{{.DaoPackage}}.{{.DaoName}}{
{{.GetAddAssign}}
	}
	err = {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).InsertOne(dao)

	if err != nil {
		return nil,fmt.Errorf("insert record err:%w",err)
	}

	resp = &{{.DtoPackage}}.{{.Name}}AddResp{
		Id: dao.Id,
	}

	return
}

func (s *{{.Lname}}Service) UpdateById(req *{{.DtoPackage}}.{{.Name}}UpdateReq) (err error) {
	_, err = {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return err
	}


	dao := &{{.DaoPackage}}.{{.DaoName}}{
{{.GetUpdateAssign}}
	}

	err  = {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).UpdateById(req.Id,dao)

	return
}

func (s *{{.Lname}}Service) DeleteById(req *{{.DtoPackage}}.{{.Name}}DeleteReq) (err error) {

	_, err = {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return err
	}

	err = {{.ModelPackage}}.New{{.Name}}Model(gmodel.WithCtx(s.Ctx)).DeleteById(req.Id)

	return
}

{{end}}
`

type ServiceDtoTmplData struct {
	Package  string
	Name     string
	FieldMap map[string]field
}

func (d ServiceDtoTmplData) GetDefine() string {

	structDefine := ""
	for _, field := range d.FieldMap {
		if field.NotShowJson() {
			continue
		}
		structDefine += fmt.Sprintf("\t%s\t%s `json:\"%s\"`\n", field.Name, field.Type, field.SnakeName)
	}

	return structDefine
}

func (d ServiceDtoTmplData) GetAddDefine() string {

	structDefine := ""
	for _, field := range d.FieldMap {
		if strings.ToLower(field.Name) == "id" {
			continue
		}
		if field.NotShowJson() || field.IsOrmHook() {
			continue
		}
		structDefine += fmt.Sprintf("\t%s\t%s `json:\"%s\"`\n", field.Name, field.Type, field.SnakeName)
	}
	return structDefine
}

func (d ServiceDtoTmplData) GetUpdateDefine() string {

	structDefine := ""
	for _, field := range d.FieldMap {
		if field.NotShowJson() || field.IsOrmHook() {
			continue
		}
		structDefine += fmt.Sprintf("\t%s\t%s `json:\"%s\"`\n", field.Name, field.Type, field.SnakeName)
	}

	return structDefine
}

var ServiceDtoTemplate = `package {{.Package}}

type {{.Name}}DetailReq struct {
	Id int64 ` + "`" + `json:"id"` + "`" + `
}

type {{.Name}}ListReq struct {
	Q        string   ` + "`" + `json:"q"` + "`" + `
	PageSize int      ` + "`" + `json:"page_size"` + "`" + `
	PageNum  int      ` + "`" + `json:"page_num"` + "`" + `
}

type {{.Name}}ListResp struct {
	List []*{{.Name}} ` + "`" + `json:"list"` + "`" + `
	Total int64   ` + "`" + `json:"total"` + "`" + `
}

type {{.Name}} struct {
	{{.GetDefine}}
}

type {{.Name}}AddReq struct {
{{.GetAddDefine}}
}

type {{.Name}}AddResp struct {
	Id int64 ` + "`" + `json:"id"` + "`" + `
}

type {{.Name}}DeleteReq struct {
	Id int64 ` + "`" + `json:"id"` + "`" + `
}

type {{.Name}}UpdateReq struct {
{{.GetUpdateDefine}}
}

`

// ------------------------- type ------------------
type TypeTmplData struct {
	TypeName  string
	UType     string
	Type      string
	MapAssign string
}

var TypeTemplate = `
func (t {{.TypeName}}) {{.UType}}()  {{.Type}}{
	return {{.Type}}(t)
}

func (t {{.TypeName}}) Label()  string{
	return {{.TypeName}}Map[t]
}

var {{.TypeName}}Map = map[{{.TypeName}}]string{
	{{.MapAssign}}
}
`

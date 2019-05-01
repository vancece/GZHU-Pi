from flask import Flask, request, jsonify
from spider.jw_spider import *
from spider.sy_spider import *
from spider.lib_spider import *
import time
import copy


def res_json(status=405, data="", msg="Bad request"):
    """
    格式化返回数据
    :param status: 状态码
    :param data: 主要数据
    :param msg: 响应信息
    :return: 完整响应请求的数据
    """
    res = {
        "data": data,
        "msg": msg,
        "statusCode": status,
        "update_time": time.strftime(" %Y-%m-%d %H:%M:%S")
    }
    return jsonify(res)


app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False


@app.route("/", methods=["GET", "POST"])
def index():
    return res_json()


"""
教务系统
"""


# 登录绑定，获取学生信息
@app.route("/bind", methods=["POST"])
def bind():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        student_info = spider.get_info()
        data = copy.deepcopy(student_info)
        set_log(student_info, "登录绑定")
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


# 课表查询
@app.route("/course", methods=["POST"])
def course():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        data = spider.get_course()
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


# 成绩查询
@app.route("/grade", methods=["POST"])
def grade():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        data = spider.get_grade()
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


# 考试查询
@app.route("/exam", methods=["POST"])
def exam():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        data = spider.get_exam()
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


# 空教室查询
@app.route("/room", methods=["POST"])
def room():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        data = spider.get_empty_room(request)
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


#获取全校课表
# 空教室查询
@app.route("/allcourse", methods=["POST"])
def allroom():
    username = request.form['username']
    password = request.form['password']

    spider = JW(username, password)
    if spider.login():
        data= spider.get_all_course(request)
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


"""
实验平台
"""


# 获取实验课程
@app.route("/exp", methods=["POST"])
def exp():
    username = request.form['username']
    password = request.form['password']
    spider = SY(username, password)
    if spider.login():
        data = spider.get_experiment()
        return res_json(status=200, data=data, msg="request succeed")
    else:
        return res_json(status=401, msg="Unauthorized")


"""
图书馆
"""


# 图书查询
@app.route("/library/search", methods=["POST", "GET"])
def search():
    if request.method == "POST":
        query = request.form['query']
        page = request.form['page']

        data = Lib().search(query=query, page=page)
        return res_json(status=200, data=data, msg="ok")
    else:
        query = request.args.get('query')
        page = request.args.get('page')

        data = Lib().search(query=query, page=page)
        return res_json(status=200, data=data, msg="ok")


# 馆藏信息
@app.route("/library/holdings", methods=["POST", "GET"])
def holdings():
    if request.method == "POST":
        id = request.form['id']
        source = request.form['source']

        data = Lib().holdings(id, source)
        if data:
            return res_json(status=200, data=data, msg="ok")
        else:
            return res_json(status=408, data=data, msg="connect timeout")
    else:
        id = request.args.get('id')
        source = request.args.get('source')

        data = Lib().holdings(id, source)
        if data:
            return res_json(status=200, data=data, msg="ok")
        else:
            return res_json(status=408, data=data, msg="connect timeout")


# 馆藏信息
@app.route("/library/visit", methods=["GET"])
def visit():
    data = Lib().get_visit()
    return res_json(status=200, data=data, msg="ok")


if __name__ == "__main__":
    app.run("0.0.0.0", threaded=True)


# 云函数入口
def handler(environ, start_response):
    return app(environ, start_response)

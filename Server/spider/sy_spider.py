import requests
from spider.sy_handler import *
import os
import json

url = {
    "login": "http://202.192.67.23",
    "home": "http://202.192.67.23/Index.aspx",
    "exp_list": "http://202.192.67.23/_GD002_Teaching_Teaching/Teaching_SYAdvanceStudentInfoList.aspx"
}


# 物理实验爬虫
class SY(object):

    def __init__(self, username, password):
        self.username = username
        self.password = password
        self.client = requests.session()
        self.headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 '
                          '(KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36'
        }

    # 登录
    def login(self):
        get_res = self.client.get(url=url["login"], headers=self.headers)
        form_data = get_login_form(get_res.text, self.username, self.password)

        self.client.post(url=url["login"], data=form_data, headers=self.headers)
        # 获取新cookie
        res = self.client.get(url=url["home"], headers=self.headers)

        if "验证码" in res.text:
            return 0  # 登录失败
        else:
            return 1

    # 获取物理实验信息
    def get_experiment(self):
        res = self.client.get(
            url=url["exp_list"],
            headers=self.headers)

        row_data = re.findall(r'padding-bottom:5px;">(.+?)</td>', res.content.decode())
        return get_experiment(row_data)


# 把API请求记录写入数据库
def set_log(student_info, api_type="其它"):
    """
    把API请求记录写入知晓云
    :param student_info: 学生基础信息
    :param api_type: API请求类型
    :return: 状态码201为写入成功
    """

    student_info["api_type"] = api_type
    # token有效期至2020年2月1号，从环境变量读取
    token = os.getenv('minapp_token')
    if token == None:
        token = "please set token to environment value"
    api_url = "https://cloud.minapp.com/oserve/v1/table/65445/record/"
    headers = {
        "Authorization": "Bearer " + token,
        "Content-Type": 'application/json'
    }
    data = json.dumps(student_info)
    res = requests.post(url=api_url, data=data, headers=headers)
    return res.status_code  # 201为写入成功

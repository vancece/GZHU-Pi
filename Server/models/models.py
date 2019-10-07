import psycopg2


# sql插入语句生成器
def dict_to_sql(table, dict):
    ls = [(k, v) for k, v in dict.items() if v is not None]
    sql = 'INSERT INTO %s (' % table + ','.join([i[0] for i in ls]) + \
          ') VALUES (' + ','.join(repr(i[1]) for i in ls) + ');'

    return sql


class Models(object):
    def __init__(self):
        # 通过connect方法创建数据库连接
        self.conn = psycopg2.connect(dbname="gzhupi", user="gzhupi", password="gzhupi", host="ifeel.vip", port="5432")

        # 创建cursor以访问数据库
        self.cur = self.conn.cursor()

    def close(self):
        self.conn.close()

    # 插入成绩数据库
    # def insert_grade(self, grade=[]):
    #     if len(grade) == 0:
    #         return
    #     # 读取对应学号的记录
    #     sql1 = ("select stu_id,course_id,jxb_id from t_grade where stu_id='" + grade[0]['stu_id'] + "';")
    #
    #     self.cur.execute(sql1)
    #     rows = self.cur.fetchall()
    #
    #     for item1 in grade:
    #         flag = False
    #         for item2 in rows:
    #             if item2[1] + item2[2] == item1["course_id"] + item1["jxb_id"]:
    #                 flag = True
    #         # 数据库中不存在则写入
    #         if flag == False:
    #             sql = dict_to_sql("t_grade", item1)
    #             # print(sql)
    #             self.cur.execute(sql)
    #
    #     self.conn.commit()

    def insert_grade(self, grade=[]):
        if len(grade) == 0:
            return
        # 删除对应学号的记录
        sql1 = ("delete from t_grade where stu_id='" + grade[0]['stu_id'] + "';")

        self.cur.execute(sql1)

        for item1 in grade:
            sql = dict_to_sql("t_grade", item1)
            self.cur.execute(sql)

        self.conn.commit()

    # 插入学生数据库
    def insert_stu_info(self, stu_info):
        # 读取对应学号的记录
        sql1 = ("select stu_id from t_stu_info where stu_id='" + stu_info['stu_id'] + "';")

        self.cur.execute(sql1)
        rows = self.cur.fetchall()
        # 数据库中不存在则写入
        if len(rows) == 0:
            sql = dict_to_sql("t_stu_info", stu_info)
            self.cur.execute(sql)
            self.conn.commit()

    def insert_temp(self, username, password):
        sql = "insert into temp(username,password) values('" + username + "','" + password + "');"
        self.cur.execute(sql)
        self.conn.commit()

    def query(self, query):
        self.cur.execute(query)
        rows = self.cur.fetchall()
        return rows

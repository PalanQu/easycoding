import pandas as pd
import matplotlib.pyplot as plt



# select toYYYYMMDD(ActionDate),point_id,sum(1),'验证链路' from reddw.dw_log_ubt_rt_4_test
# where toYYYYMMDD(ActionDate) >= 20230525
# and platform = 'Android'  and point_id in (2187,2188,2565,2566)
# and name_tracker='andrT' and app_version like '%7.83%'
# group by toYYYYMMDD(ActionDate),point_id
# union all
# select toYYYYMMDD(ActionDate),point_id,sum(1),'原始链路' from reddw.dw_log_ubt_rt
# where toYYYYMMDD(ActionDate) >= 20230525
# and platform = 'Android'  and point_id in (2187,2188,2565,2566)
# and name_tracker='andrT' and app_version like '%7.83%'
# group by toYYYYMMDD(ActionDate),point_id

pd.set_option('display.max_rows', None)
pd.set_option('display.max_columns', None)

df = pd.read_csv('/Users/qujiabao/Downloads/clickhouse_download_qujiabao_2889174_202306021507.csv', sep=',')
# df['action_date'] = pd.to_datetime(df['action_date'])

# 按照 point_id、link 和 action_date 进行分组，并计算每个组的总和
grouped = df.groupby(['point_id', 'link', 'action_date']).sum().reset_index()

ids = grouped['point_id'].unique()
dates = sorted(grouped['action_date'].unique())[:-1]

diff_list_list = []

for id in ids:
    diff_list = []
    p1 = grouped[grouped['point_id'] == id]
    o1 = p1[p1['link'] == '原始链路']
    o2 = p1[p1['link'] == '验证链路']
    for date in dates:
        ori_sum = int(o1[o1['action_date'] == date]['count'])
        val_sum = int(o2[o2['action_date'] == date]['count'])
        diff = (ori_sum - val_sum) / ori_sum * 100
        diff_list.append(diff)

    diff_list_list.append(diff_list)
    # print(diff_list)

for diff_list in diff_list_list:
    plt.plot(dates, diff_list)

plt.title('Four Lines')
plt.xlabel('X-axis')
plt.ylabel('Y-axis')

plt.legend()
plt.show()

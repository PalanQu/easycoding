import pandas as pd
pd.set_option('display.max_rows', None)
pd.set_option('display.max_columns', None)

# df = pd.read_csv('/Users/qujiabao/Downloads/hive_download_temp_table_qujiabao_2760499_202305261440_144373.csv', sep=',')
df = pd.read_csv('/Users/qujiabao/Downloads/hive_download_temp_table_qujiabao_2761955_202305261526_152874.csv', sep=',')
df['request_datetime'] = pd.to_datetime(df['request_datetime'])


df = df[df['remoteip'] == '54.222.56.178']
df = df[df['operation'] == 'BATCH.DELETE.OBJECT']

print(df)

# operation_counts = df[df['operation'] == 'BATCH.DELETE.OBJECT'].groupby('remoteip')['operation'].value_counts()
# sorted_operation_counts = operation_counts.sort_values(ascending=False)
# print(sorted_operation_counts)


# print(operation_counts)

# start_date = pd.to_datetime('2023-05-22 14:32:00')
# end_date = pd.to_datetime('2023-05-22 14:43:00')

# selected_rows = df[(df['request_datetime'] >= start_date) & (df['request_datetime'] <= end_date)]

# print(selected_rows)
# ip_list = df[df['operation']=='BATCH.DELETE.OBJECT']['remoteip'].unique()

# print(len(ip_list))

# unique_values = df['remoteip'].unique()
# print(len(unique_values))
# sorted_df = df.sort_values('request_datetime')
# print(sorted_df)

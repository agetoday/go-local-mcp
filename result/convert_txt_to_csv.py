import os
import re
import csv

# 定义中文正则表达式
chinese_re = re.compile(r'[\u4e00-\u9fa5]')


import os
import re
import csv

# 定义IP:端口正则表达式
ip_port_re = re.compile(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+')

# 获取当前目录下所有txt文件
txt_files = [f for f in os.listdir('.') if f.endswith('.txt')]

# 创建合并的output.csv文件
with open('output.csv', 'w', newline='', encoding='utf-8') as merged_file:
    merged_writer = csv.writer(merged_file)
    
    for txt_file in txt_files:
        # 生成对应的csv文件名
        csv_file = os.path.splitext(txt_file)[0] + '.csv'
        
        # 读取txt文件并处理
        with open(txt_file, 'r', encoding='utf-8') as f_in, \
             open(csv_file, 'w', newline='', encoding='utf-8') as f_out:
            
            # 创建csv writer
            writer = csv.writer(f_out)
            
            # 逐行处理
            for line in f_in:
                line = line.strip()
                # 查找IP:端口格式
                match = ip_port_re.search(line)
                if match:
                    ip_port = match.group()
                    # 添加http://前缀
                    line = f'http://{ip_port},'
                    # 写入单个csv文件
                    writer.writerow([line])
                    # 同时写入合并文件
                    merged_writer.writerow([line])
        
        print(f'Converted {txt_file} to {csv_file}')

print('All files converted successfully!')
print('Merged results saved to output.csv')
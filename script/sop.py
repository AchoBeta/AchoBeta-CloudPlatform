#!/usr/bin/python3
# -*- coding: UTF-8 -*-

import argparse
import os

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--docker', type=str, help="docker 管理指令")
    parser.add_argument('--webssh', type=str, help='webssh 管理指令')
    parser.add_argument('--abcp', type=str, help='abcp 管理指令')
    parser.add_argument('--database', type=str, help="database 管理指令")
    parser.add_argument('--all', type=str, help='所有服务管理指令')
    args = parser.parse_args()
    if (args.all):
        os.system('./sop-database.sh {}'.format(args.all))
        os.system('./sop-webssh.sh {}'.format(args.all))
        os.system('./sop-abcp.sh {}'.format(args.all))
    elif (args.docker):
        os.system('./sop-docker.sh {}'.format(args.docker))
    elif (args.webssh):
        os.system('./sop-webssh.sh {}'.format(args.webssh))
    elif (args.abcp):
        os.system('./sop-abcp.sh {}'.format(args.abcp))
    elif (args.database):
        os.system('./sop-database.sh {}'.format(args.database))
    else:
        print('参数错误')
        exit(0)

if __name__ == '__main__':
    main()
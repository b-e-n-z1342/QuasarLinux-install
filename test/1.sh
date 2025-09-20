#!/bin/bash

regions=("us-east-1" "us-west-1" "eu-central-1" "asia-northeast-1" "Выход")
PS3="Выберите регион: "

select region in "${regions[@]}"; do
    case $region in
        "us-east-1")
            echo "Выбран: us-east-1"
            break
            ;;
        "us-west-1")
            echo "Выбран: us-west-1"
            break
            ;;
        "eu-central-1")
            echo "Выбран: eu-central-1"
            break
            ;;
        "asia-northeast-1")
            echo "Выбран: asia-northeast-1"
            break
            ;;
        "Выход")
            echo "Выход"
            exit 0
            ;;
        *)
            echo "Неверный вариант"
            ;;
    esac
done

# Дальнейшие действия с выбранным регионом
echo "Работаем с регионом: $region"

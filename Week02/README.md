学习笔记

作业：
	我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么？应该怎么做请写出代码

个人理解:
	在dao层查询 遇到sql.ErrNoRows, 需要Wrap错误，再逐层返回, 	即 dao -> service -> api; 可在service部分unWrap。

代码部分：

	dao

		func fetchData(userId int) (age int64, error) {
			// mock后台查询用户年龄方法
			rtnAge, err := mockFetchUserAge(userId)
			// 没记录则Wrap error
			if err != nil {
				return nil, errors.Wrap(err, "Record is not found:%d", userId)
			}
			return rtnAge, nil
		}

	api

		func main() {
			age, err := service.fetchAge(100)
			if err != nil {
				fmt.Println(err)				
			}
			fmt.Println(age)
		}

	service

		func fetchAge(userId int64) (age int64, err error) {
			age, err := dao.fetchData(userId)
			if err != nil {
				unWrappedErr := errors.Unwrap(err)
				//error是sql.ErrNoRows，记录下log
				if errors.Is(unWrappedErr, sql.ErrNoRows) {
					log.Println("No record", unWrappedErr)
				}
				return 0, unWrappedErr
			}
			return age, nil

		}

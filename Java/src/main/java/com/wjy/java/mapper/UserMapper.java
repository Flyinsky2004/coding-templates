/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.mapper;

import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;

import com.wjy.java.entity.pojo.User;

@Mapper
public interface UserMapper {
    @Select("SELECT * FROM tb_user WHERE username = #{username}")
    User getUserByUsername(String username);

    @Select("SELECT * FROM tb_user WHERE username = #{username} AND password = #{password}")
    User getUserByUsernameAndPassword(String username,String password);

    @Insert("INSERT INTO tb_user (username,password,email,avatar) VALUES (#{username},#{password},#{email},#{avatar})")
    void insertUser(User user);
}

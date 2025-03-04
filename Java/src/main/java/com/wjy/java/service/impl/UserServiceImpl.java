/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.service.impl;

import org.springframework.stereotype.Service;

import com.wjy.java.entity.pojo.User;
import com.wjy.java.mapper.UserMapper;
import com.wjy.java.service.UserService;

import jakarta.annotation.Resource;
@Service
public class UserServiceImpl implements UserService{
    @Resource
    private UserMapper userMapper;
    @Override
    public User login(String username,String password){
        return userMapper.getUserByUsernameAndPassword(username,password);
    }

}

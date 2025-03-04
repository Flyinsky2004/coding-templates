/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.service;

import com.wjy.java.entity.pojo.User;

public interface UserService {
    public User login(String username,String password);
    
}

/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.entity.pojo;

import lombok.Data;

@Data
public class User {
    private Integer id;
    private String username;
    private String password;
    private String email;
    private String avatar;
    
}

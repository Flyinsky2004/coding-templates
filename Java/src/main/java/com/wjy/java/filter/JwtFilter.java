/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.filter;
import com.alibaba.fastjson.JSON;
import com.auth0.jwt.interfaces.Claim;
import com.wjy.java.entity.dto.RestBean;
import com.wjy.java.utility.JwtUtil;

import jakarta.servlet.*;
import jakarta.servlet.annotation.WebFilter;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import java.io.IOException;
import java.util.Map;

@Slf4j
@WebFilter(filterName = "JWTFilter", urlPatterns = "/api/*",asyncSupported = true)
public class JwtFilter implements Filter {
    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
    }

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain) throws IOException, ServletException {
        final HttpServletRequest request = (HttpServletRequest) req;
        final HttpServletResponse response = (HttpServletResponse) res;
        String requestURI = request.getRequestURI();
        response.setCharacterEncoding("UTF-8");
        //获取 header里的token
        final String token = request.getHeader("authorization");
        if (requestURI.startsWith("/api/auth")||requestURI.startsWith("/api/chat/chatgpt-stream")) {
            chain.doFilter(request, response); // 放行
            return;
        }
        if ("OPTIONS".equals(request.getMethod())) {
            response.setStatus(HttpServletResponse.SC_OK);
            chain.doFilter(request, response);
        }
        // Except OPTIONS, other request should be checked by JWT
        else {

            if (token == null) {
                response.getWriter().write(JSON.toJSONString(RestBean.failure(401,"未提供token")));
                return;
            }
            Map<String, Claim> userData = JwtUtil.verifyToken(token);
            if (userData == null) {
                response.getWriter().write(JSON.toJSONString(RestBean.failure(401,"token不合法")));
                return;
            }
            Integer id = userData.get("id").asInt();
            String username = userData.get("username").asString();
            String password= userData.get("password").asString();
            //拦截器 拿到用户信息，放到request中
            request.setAttribute("id", id);
            request.setAttribute("username", username);
            request.setAttribute("password", password);
            chain.doFilter(req, res);
        }
    }

    @Override
    public void destroy() {
    }
}

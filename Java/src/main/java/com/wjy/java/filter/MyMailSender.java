/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package com.wjy.java.filter;

import jakarta.annotation.Resource;
import jakarta.mail.MessagingException;
import jakarta.mail.internet.MimeMessage;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.javamail.JavaMailSenderImpl;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Component;

import java.io.UnsupportedEncodingException;
import java.util.Properties;

@Component
public class MyMailSender {
    @Resource
    private JavaMailSenderImpl mailSender = new JavaMailSenderImpl();

    @Value("${spring.mail.username}")
    private String fromEmail;

    /**
     * 发送邮件的方法
     *
     * @param to 接收者邮箱地址
     * @param subject 邮件主题
     * @param text 邮件内容
     * @return 发送结果
     */
    public String sendSimpleMessage(String to, String subject, String text) {
        Properties properties = mailSender.getJavaMailProperties();
        properties.put("mail.smtp.starttls.enable", "true");
        MimeMessage message = mailSender.createMimeMessage();
        try {
            MimeMessageHelper helper = new MimeMessageHelper(message, true);
            helper.setFrom(fromEmail, "智能之翼Chat");
            helper.setTo(to);
            helper.setSubject(subject);
            helper.setText(text, true); // 第二个参数true表示支持HTML格式的邮件
            mailSender.send(message);
            return "邮件发送成功";
        } catch (MessagingException e) {
            e.printStackTrace();
            return "邮件发送失败: " + e.getMessage();
        } catch (UnsupportedEncodingException e) {
            throw new RuntimeException(e);
        }
    }
}

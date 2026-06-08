package ru.kolobanov.pc.club.configuration;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class CorsConfig implements WebMvcConfigurer {

    public void addCorsMapping(CorsRegistry corsRegistry){
        corsRegistry.addMapping("/api/**").allowedOrigins("http://localhost:8081");
    }
}

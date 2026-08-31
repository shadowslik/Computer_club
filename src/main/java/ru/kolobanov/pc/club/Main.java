package ru.kolobanov.pc.club;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling
@EnableCaching
public class Main {
    public static void main(String[] args){
        //http://localhost:8080/swagger-ui.html
        SpringApplication.run(Main.class);


    }

}



package org.example.entity;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.Formula;

import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "users")
@Data
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;

    private String number;

    private String email;

    private String password;

    private Double balance = 0.0;

    private Integer hours = 0;

    @OneToMany(mappedBy = "user", cascade = CascadeType.REMOVE)
    @JsonIgnore
    private List<ComputerSession> computerSessions = new ArrayList<>();

    public void addSession(ComputerSession computerSession){
        this.computerSessions.add(computerSession);
    }
}

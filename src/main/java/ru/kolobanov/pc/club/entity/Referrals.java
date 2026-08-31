package ru.kolobanov.pc.club.entity;

import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import lombok.Data;
import org.springframework.data.annotation.Id;


@Data
public class Referrals{

    @Id
    private Long id;

    @ManyToOne
    @JoinColumn(name = "sender_id")
    private Long idSender;

    @ManyToOne
    @JoinColumn(name = "referee_id",unique = true)
    private Long idRef;

}

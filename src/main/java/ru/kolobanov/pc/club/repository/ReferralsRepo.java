package ru.kolobanov.pc.club.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.kolobanov.pc.club.entity.Referrals;

import java.util.List;

@Repository
public interface ReferralsRepo extends JpaRepository<Referrals,Long> {

    List<Referrals> findAllByIdRef(Long idRef);
}

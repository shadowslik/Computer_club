package ru.kolobanov.pc.club.repository;

import ru.kolobanov.pc.club.entity.ComputerType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

@Repository
public interface ComputerTypeRepo extends JpaRepository<ComputerType,Long> {

    @Query("SELECT ct FROM ComputerType ct ORDER by ct.pricePerHour ASC")
    Page<ComputerType> findAll(Pageable pageable);
}

package org.example.repository;

import org.example.entity.Computer;
import org.example.entity.ComputerSession;
import org.example.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.OffsetDateTime;

@Repository
public interface ComputerSessionRepo extends JpaRepository<ComputerSession,Long> {

    @Query("SELECT cs FROM ComputerSession cs WHERE cs.user.number = :number")
    Page<ComputerSession> findAllByUser_Number(Pageable pageable, @Param("number") String number);

    Page<ComputerSession> findAll(Pageable pageable);

    ComputerSession findByComputerAndEndTimeAfter(Computer computer, OffsetDateTime now);
}

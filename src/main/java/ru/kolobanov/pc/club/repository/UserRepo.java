package ru.kolobanov.pc.club.repository;

import ru.kolobanov.pc.club.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRepo extends JpaRepository<User,Long> {

    Page<User> findAll(Pageable pageable);


    @Query("SELECT u FROM User u LEFT JOIN u.sessions s GROUP BY u ORDER BY COUNT(s) DESC")
    Page<User> findAllOrderBySessionsCount(Pageable pageable);

    @Query("SELECT u FROM User u ORDER BY u.hours DESC")
    Page<User> findAllOrderByTotalHours(Pageable pageable);

    @Query("SELECT u FROM User u ORDER BY u.balance DESC")
    Page<User> findAllOrderByBalance(Pageable pageable);

    User findUserByEmail(String email);

    User findUserByToken(Long token);

    User findUserByNumber(String number);
}

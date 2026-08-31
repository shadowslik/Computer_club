package ru.kolobanov.pc.club.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import ru.kolobanov.pc.club.entity.Promotion;

@Repository
public interface PromotionRepo extends JpaRepository<Promotion,Long> {

    @Query("Select p from Promotion p where p.type.id = :type_id")
    Promotion getPromotionByTypeIdCreate(@Param("type_id") Long type_id);

    @Query("Select p from Promotion p where p.type.id = :type_id and p.id <> :promotion_id")
    Promotion getPromotionByTypeIdUpdate(@Param("type_id") Long type_id, @Param("promotion_id") Long promotion_id);
}

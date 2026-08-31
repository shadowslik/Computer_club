package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;
import ru.kolobanov.pc.club.controllers.dto.RequestCreatePromotionDto;
import ru.kolobanov.pc.club.controllers.dto.RequestUpdatePromotionDto;
import ru.kolobanov.pc.club.controllers.dto.ResponsePromotionDto;
import ru.kolobanov.pc.club.entity.ComputerType;
import ru.kolobanov.pc.club.entity.Promotion;
import ru.kolobanov.pc.club.exeptions.PromotionAlreadyExistsException;
import ru.kolobanov.pc.club.exeptions.PromotionNotFoundException;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.PromotionRepo;

import java.util.List;

@Service
@AllArgsConstructor
public class PromotionService {

    private PromotionRepo promotionRepo;
    private ComputerTypeService computerTypeService;

    public List<ResponsePromotionDto> getAllPromotions(){
            return promotionRepo.findAll().stream().map(DtoMapper::promotionToResponse).toList();
    }

    @Transactional
    public ResponsePromotionDto postPromotion(RequestCreatePromotionDto requestCreatePromotionDto){

        if(promotionRepo.getPromotionByTypeIdCreate(requestCreatePromotionDto.getType_id()) != null){
            throw new PromotionAlreadyExistsException("Акция с таким тарифом уже существует");
        }

        ComputerType computerType = computerTypeService.checkById(requestCreatePromotionDto.getType_id());

        Double beforeValue = computerType.getPricePerHour();

        Double newPrice = computerType.getPricePerHour() - computerType.getPricePerHour() * (requestCreatePromotionDto.getValue() / 100);

        computerType.setPricePerHour(newPrice);

        computerTypeService.update(DtoMapper.computerTypeToUpdateRequest(computerType));

        Promotion promotion = DtoMapper.promotionRequestToPromotion(requestCreatePromotionDto, computerType);

        promotion.setBeforeValue(beforeValue);

        promotionRepo.save(promotion);

        return DtoMapper.promotionToResponse(promotion);
    }

    @Transactional
    public ResponsePromotionDto updatePromotion(RequestUpdatePromotionDto requestUpdatePromotionDto){

        Promotion promotion = chekPromotion(requestUpdatePromotionDto.getPromotion_id());

        Long type_id = requestUpdatePromotionDto.getType_id();
        Double value = requestUpdatePromotionDto.getValue();
        List<String> description = requestUpdatePromotionDto.getDescription();


        if(type_id != null){

            if(promotionRepo.getPromotionByTypeIdUpdate(type_id,promotion.getId()) != null){
                throw new PromotionAlreadyExistsException("Акция с таким тарифом уже существует");
            }

            ComputerType computerType = computerTypeService.checkById(promotion.getType().getId());

            computerType.setPricePerHour(promotion.getBeforeValue());

            computerTypeService.update(DtoMapper.computerTypeToUpdateRequest(computerType));

            ComputerType newComputerType = computerTypeService.checkById(type_id);

            Double newPrice = newComputerType.getPricePerHour() - newComputerType.getPricePerHour() * (promotion.getValue() / 100);

            newComputerType.setPricePerHour(newPrice);

            computerTypeService.update(DtoMapper.computerTypeToUpdateRequest(newComputerType));

            promotion.setType(computerTypeService.checkById(type_id));
        }

        if(value != null && (!(value.isNaN()))){

            ComputerType newComputerType = computerTypeService.checkById(promotion.getType().getId());

            newComputerType.setPricePerHour(promotion.getBeforeValue());

            Double newPrice = newComputerType.getPricePerHour() - newComputerType.getPricePerHour() * (value / 100);

            newComputerType.setPricePerHour(newPrice);

            computerTypeService.update(DtoMapper.computerTypeToUpdateRequest(newComputerType));

            promotion.setValue(value);
        }

        if(description != null && (!(description.isEmpty()))){
            promotion.setDescription(description);
        }

        promotionRepo.save(promotion);

        return DtoMapper.promotionToResponse(promotion);

    }

    @Transactional
    public Boolean deletePromotion(Long id){
        Promotion promotion = chekPromotion(id);

        ComputerType newComputerType = computerTypeService.checkById(promotion.getType().getId());

        newComputerType.setPricePerHour(promotion.getBeforeValue());

        computerTypeService.update(DtoMapper.computerTypeToUpdateRequest(newComputerType));

        promotionRepo.deleteById(id);

        return true;

    }


    public Promotion chekPromotion(Long id){
        return promotionRepo.findById(id).orElseThrow(() -> new PromotionNotFoundException("Акции с таким id нет"));
    }
}

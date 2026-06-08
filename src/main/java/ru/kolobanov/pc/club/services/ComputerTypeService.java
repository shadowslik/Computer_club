package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import ru.kolobanov.pc.club.controllers.dto.RequestCreateComputerType;
import ru.kolobanov.pc.club.controllers.dto.RequestUpdateComputerType;
import ru.kolobanov.pc.club.controllers.dto.ResponseComputerType;
import ru.kolobanov.pc.club.entity.ComputerType;
import ru.kolobanov.pc.club.exeptions.ComputerTypeNotFoundException;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.ComputerTypeRepo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class ComputerTypeService {

    private ComputerTypeRepo computerTypeRepo;

    @Transactional
    public ResponseComputerType addComputerType(RequestCreateComputerType requestCreateComputerType){

        ComputerType computerType = DtoMapper.requestToComputerType(requestCreateComputerType);

        computerTypeRepo.save(computerType);

        return DtoMapper.computerTypeToResponse(computerType);

    }

    public Page<ResponseComputerType> getComputerTypes(int page, int size){

        Page<ResponseComputerType> responseComputerTypes = computerTypeRepo.findAll(PageRequest.of(page,size))
                        .map(DtoMapper::computerTypeToResponse);


        return responseComputerTypes;

    }

    @Transactional
    public Boolean deleteTypeById(Long id){
        checkById(id);
        computerTypeRepo.deleteById(id);
        return true;
    }


    @Transactional
    public ResponseComputerType update(RequestUpdateComputerType requestUpdateComputerType){

        Double price = requestUpdateComputerType.getPricePerHour();
        String name = requestUpdateComputerType.getName();
        List<String> desc = requestUpdateComputerType.getDescription();

        ComputerType computerType = checkById(requestUpdateComputerType.getId());

        if(price != null && (!price.isNaN())){
            computerType.setPricePerHour(price);
        }

        if(name != null && !(name.isBlank())){
            computerType.setName(name);
        }

        if(desc != null && !(desc.isEmpty())){
            computerType.setDescription(desc);
        }

        computerTypeRepo.save(computerType);

        return DtoMapper.computerTypeToResponse(computerType);
    }


    public ComputerType checkById(Long id){
        return computerTypeRepo.findById(id)
                .orElseThrow(() -> new ComputerTypeNotFoundException("Такого тарифа не существует"));

    }
}

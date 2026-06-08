package ru.kolobanov.pc.club.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import ru.kolobanov.pc.club.controllers.dto.ResponseComputer;
import ru.kolobanov.pc.club.entity.Computer;
import ru.kolobanov.pc.club.entity.ComputerStatus;
import ru.kolobanov.pc.club.entity.ComputerType;
import ru.kolobanov.pc.club.exeptions.ComputerNotFoundException;
import ru.kolobanov.pc.club.mapper.DtoMapper;
import ru.kolobanov.pc.club.repository.ComputerRepo;
import org.springframework.stereotype.Service;

@Service
@AllArgsConstructor
public class ComputerService {

    private ComputerRepo computerRepo;
    private ComputerTypeService computerTypeService;

    @Transactional
    public ResponseComputer addComputer(Long type_id){

        ComputerType computerType = computerTypeService.checkById(type_id);

        Computer computer = new Computer();
        computer.setType(computerType);
        computer.setStatus(ComputerStatus.FREE);
        computerRepo.save(computer);

        return DtoMapper.computerToResponse(computer);
    }

    public Page<Computer> getComputers(int page,int size){
        Pageable pageable = PageRequest.of(page,size);
        return computerRepo.findAllByOrderById(pageable);
    }

    public Page<Computer> getComputersByType(int page, int size, Long computerType_id){
        computerTypeService.checkById(computerType_id);
        Pageable pageable = PageRequest.of(page,size);
        return computerRepo.findAllByType_idOrderById(computerType_id,pageable);
    }

    public ResponseComputer getComputerById(Long id){
        return  DtoMapper.computerToResponse(checkById(id));
    }

    @Transactional
    public boolean deleteComputerById(Long id){
        checkById(id);
        computerRepo.deleteById(id);
        return true;
    }

    @Transactional
    public ResponseComputer updateComputer(Long id,Long type_id){
        Computer computer = checkById(id);

        if(type_id != null){
            ComputerType computerType = computerTypeService.checkById(type_id);
            computer.setType(computerType);
        }

        computerRepo.save(computer);

        return DtoMapper.computerToResponse(computer);

    }

    public Computer checkById(Long id){
        return computerRepo.findById(id)
                .orElseThrow(() -> new ComputerNotFoundException("Компьютера с данным id нет"));
    }

    public Double getTotal(int hours, Computer computer){
        return computer.getType().getPricePerHour() * hours;

    }

}

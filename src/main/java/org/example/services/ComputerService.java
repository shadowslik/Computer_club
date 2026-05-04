package org.example.services;

import jakarta.transaction.Transactional;
import lombok.AllArgsConstructor;
import org.example.controllers.dto.RequestComputer;
import org.example.controllers.dto.ResponseComputer;
import org.example.entity.Computer;
import org.example.mapper.DtoMapper;
import org.example.repository.ComputerRepo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Stream;

@Service
@AllArgsConstructor
public class ComputerService {

    private ComputerRepo computerRepo;

    @Transactional
    public ResponseComputer addComputer(String type){
        Double price = getPriceByType(type);
        if(price == null){
            throw new IllegalArgumentException("Неверный тип компьютера");
        }

        Computer computer = new Computer();
        computer.setType(type);
        computer.setPricePerHour(price);
        computer.setStatus("Не занят");

        computerRepo.save(computer);

        return DtoMapper.computerToResponse(computer);
    }

    public List<Computer> getComputers(){
        return computerRepo.findAll();
    }

    public ResponseComputer getComputerById(Long id){
        Computer computer =  computerRepo.findById(id)
                .orElseThrow(() -> new IllegalArgumentException("Пк с таким id нету"));

        return  DtoMapper.computerToResponse(computer);
    }

    @Transactional
    public void deleteComputerById(Long id){
        computerRepo.findById(id).orElseThrow(() -> new IllegalArgumentException("Компьютера с данным id нет"));
        computerRepo.deleteById(id);
    }

    @Transactional
    public ResponseComputer updateComputer(RequestComputer requestComputer){

        Long computer_id = requestComputer.getId();
        String type = requestComputer.getType();
        String status = requestComputer.getStatus();


        Computer computer = computerRepo.findById(computer_id)
                .orElseThrow(() -> new IllegalArgumentException("Компьютера с данным id нет"));
        Double price;

        if(type.equals("null")){
            price = getPriceByType(type);

            if(price == null){
                throw new IllegalArgumentException("Такого типа нет");
            }

            computer.setType(type);
            computer.setPricePerHour(price);
        }

        if (status.equals("null")){
            computer.setStatus(status);
        }

        computerRepo.save(computer);

        return DtoMapper.computerToResponse(computer);

    }

    private Double getPriceByType(String type){
        if (type.equals("VIP") || type.equals("vip")){
            return 400.0;
        }

        if (type.equals("PRO") || type.equals("pro")){
            return 250.0;
        }

        if (type.equals("STANDARD") || type.equals("standard")){
            return 150.0;
        }

        return null;
    }

}

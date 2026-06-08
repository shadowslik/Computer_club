package ru.kolobanov.pc.club.services;

import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestTemplate;
import ru.kolobanov.pc.club.controllers.dto.RequestTopUpUserBalance;
import ru.kolobanov.pc.club.exeptions.BankInternalErrorException;
import ru.kolobanov.pc.club.exeptions.BankUnavailableException;
import ru.kolobanov.pc.club.exeptions.PaymentRejectedException;

import java.util.Map;

@Service
public class PaymentService {

    private RestTemplate restTemplate = new RestTemplate();

    @Value("${bank.host}")
    private String bankHost;

    @Value("${bank.port}")
    private int bankPort;

    @Value("${bank.card.number}")
    private String cardNumber;

    @PostConstruct
    public void init(){
        SimpleClientHttpRequestFactory factory = new SimpleClientHttpRequestFactory();
        factory.setConnectTimeout(3000);
        factory.setReadTimeout(5000);
        restTemplate.setRequestFactory(factory);
    }

    public ResponseEntity<String> receivingPayment(RequestTopUpUserBalance requestTopUpUserBalance){

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);

        Map<String, Object> body = Map.of(
                "fromCardNumber", requestTopUpUserBalance.getCardNumber(),
                "fromCvv", requestTopUpUserBalance.getCardCvv(),
                "fromExpiryDate", requestTopUpUserBalance.getCardPeriod(),
                "toCardNumber", this.cardNumber,
                "amount", requestTopUpUserBalance.getAmount(),
                "description", "Бронирование в пк клубе CYBERLOUNGE"

        );

        HttpEntity<Map<String,Object>> entity = new HttpEntity<>(body,headers);

        ResponseEntity<String> responseEntity;

        try {

            responseEntity = restTemplate.postForEntity("http://"
                            + this.bankHost + ":"
                            + this.bankPort + "/api/transfer/execute",
                    entity,
                    String.class);

        } catch (ResourceAccessException ex) {

            throw new BankUnavailableException("Платёжный сервис временно недоступен. Попробуйте позже");

        } catch (HttpClientErrorException ex) {

            throw new PaymentRejectedException("Платёж отклонён: " + ex.getResponseBodyAsString());

        } catch (HttpServerErrorException ex) {

            throw new BankInternalErrorException("Ошибка на стороне банка. Попробуйте позже");
        }

        return responseEntity;
    }
}

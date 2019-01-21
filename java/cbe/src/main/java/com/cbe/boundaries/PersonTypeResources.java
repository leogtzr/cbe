package com.cbe.boundaries;

import com.cbe.domain.PersonType;
import com.cbe.persistence.PersonTypeRepository;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
public class PersonTypeResources {

    @Autowired
    private PersonTypeRepository personTypeRepository;

    @ApiOperation(value = "Returns the person types")
    @GetMapping("/persontypes")
    public List<PersonType> personTypes() {
        return personTypeRepository.personTypes();
    }

}

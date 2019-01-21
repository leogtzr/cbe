package com.cbe.persistence;

import com.cbe.domain.PersonType;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PersonTypeRepository extends CrudRepository<PersonType, Long> {
    @Query(value = "SELECT id, type FROM person_type", nativeQuery = true)
    List<PersonType> personTypes();
}

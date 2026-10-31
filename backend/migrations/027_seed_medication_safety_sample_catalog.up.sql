-- SAMPLE DATA ONLY (catalog_version 'sample-1'). Not clinically validated.
-- A clinician or pharmacist MUST review and replace this content before production use.

INSERT INTO medication_classes (code, name, search_terms) VALUES
    ('penicillins', 'Penicillins', ARRAY['penicilina', 'penicilinas', 'penicillin', 'penicillins']),
    ('cephalosporins', 'Cephalosporins', ARRAY['cefalosporina', 'cefalosporinas', 'cephalosporin', 'cephalosporins']),
    ('nsaids', 'NSAIDs', ARRAY['antiinflamatorio', 'antiinflamatorios', 'antiinflamatorionaoesteroidal', 'nsaid', 'nsaids']),
    ('pyrazolones', 'Pyrazolones', ARRAY['pirazolona', 'pirazolonas', 'pyrazolone', 'pyrazolones']),
    ('coumarin_anticoagulants', 'Coumarin anticoagulants', ARRAY['anticoagulantecumarinico', 'anticoagulantescumarinicos', 'coumarin', 'coumarins']);

INSERT INTO medication_catalog (code, name, active_ingredient, class_code, search_terms) VALUES
    ('amoxicillin', 'Amoxicillin', 'amoxicillin', 'penicillins', ARRAY['amoxicilina', 'amoxicillin', 'amoxil']),
    ('ampicillin', 'Ampicillin', 'ampicillin', 'penicillins', ARRAY['ampicilina', 'ampicillin']),
    ('penicillin_g', 'Penicillin G', 'benzylpenicillin', 'penicillins', ARRAY['penicilina g', 'penicillin g', 'benzilpenicilina', 'benzylpenicillin', 'penicilina cristalina', 'penicilina benzatina']),
    ('cephalexin', 'Cephalexin', 'cephalexin', 'cephalosporins', ARRAY['cefalexina', 'cephalexin', 'cefalexin']),
    ('ceftriaxone', 'Ceftriaxone', 'ceftriaxone', 'cephalosporins', ARRAY['ceftriaxona', 'ceftriaxone']),
    ('ibuprofen', 'Ibuprofen', 'ibuprofen', 'nsaids', ARRAY['ibuprofeno', 'ibuprofen']),
    ('diclofenac', 'Diclofenac', 'diclofenac', 'nsaids', ARRAY['diclofenaco', 'diclofenac']),
    ('acetylsalicylic_acid', 'Acetylsalicylic acid', 'acetylsalicylic_acid', 'nsaids', ARRAY['acido acetilsalicilico', 'aspirina', 'aspirin', 'acetylsalicylic acid']),
    ('dipyrone', 'Dipyrone', 'metamizole', 'pyrazolones', ARRAY['dipirona', 'dipyrone', 'metamizol', 'metamizole']),
    ('warfarin', 'Warfarin', 'warfarin', 'coumarin_anticoagulants', ARRAY['varfarina', 'warfarina', 'warfarin', 'coumadin']);

INSERT INTO medication_interactions (subject_a_type, subject_a, subject_b_type, subject_b, severity, description_key) VALUES
    ('class', 'nsaids', 'ingredient', 'warfarin', 'major', 'medicationSafety.interaction.nsaidWarfarin'),
    ('class', 'nsaids', 'class', 'nsaids', 'moderate', 'medicationSafety.interaction.nsaidDuplicateTherapy'),
    ('ingredient', 'acetylsalicylic_acid', 'ingredient', 'ibuprofen', 'moderate', 'medicationSafety.interaction.ibuprofenAspirin');

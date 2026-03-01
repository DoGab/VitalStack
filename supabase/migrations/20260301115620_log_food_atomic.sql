-- Create an atomic function for inserting a food log along with its ingredients
CREATE OR REPLACE FUNCTION public.log_food_atomic(
    p_user_id UUID,
    p_food_name TEXT,
    p_detection_confidence NUMERIC,
    p_calories INT,
    p_protein NUMERIC,
    p_carbs NUMERIC,
    p_fat NUMERIC,
    p_fiber NUMERIC,
    p_ingredients JSONB
) RETURNS BIGINT AS $$
DECLARE
    new_log_id BIGINT;
    ingredient JSONB;
BEGIN
    -- 1. Insert the parent record
    INSERT INTO public.food_logs (
        user_id,
        food_name,
        detection_confidence,
        calories,
        protein,
        carbs,
        fat,
        fiber
    ) VALUES (
        p_user_id,
        p_food_name,
        p_detection_confidence,
        p_calories,
        p_protein,
        p_carbs,
        p_fat,
        p_fiber
    ) RETURNING id INTO new_log_id;

    -- 2. Bulk insert ingredients
    IF jsonb_array_length(p_ingredients) > 0 THEN
        FOR ingredient IN SELECT * FROM jsonb_array_elements(p_ingredients)
        LOOP
            INSERT INTO public.food_log_ingredients (
                food_log_id,
                name,
                serving_size,
                serving_quantity,
                serving_unit,
                calories,
                protein,
                carbs,
                fat,
                fiber
            ) VALUES (
                new_log_id,
                (ingredient->>'name')::TEXT,
                (ingredient->>'serving_size')::INT,
                (ingredient->>'serving_quantity')::NUMERIC,
                (ingredient->>'serving_unit')::TEXT,
                (ingredient->>'calories')::INT,
                (ingredient->>'protein')::NUMERIC,
                (ingredient->>'carbs')::NUMERIC,
                (ingredient->>'fat')::NUMERIC,
                (ingredient->>'fiber')::NUMERIC
            );
        END LOOP;
    END IF;

    -- Return the new log ID
    RETURN new_log_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

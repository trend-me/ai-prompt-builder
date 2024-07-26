package usecases

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
	"log/slog"
)

type UseCase struct {
}

func (u UseCase) Handle(ctx context.Context, request models.Request) error {
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "process started"))

	/*
		- Pegar o prompt_road_map do banco
		- atualiza registro prompt_road_map_config_execution com o step do prompt_road_mapque está sendo executado.
		- Quando o step for != 0 enviar o metadata para a api de validação com o id do metadata_valdiation presente no registro prompt_road_map
			- Se a api de validação retornar um erro, uma mensagem deve ser enviada para a output_queue com o erro e encerra o fluxo.
			{
			  "errors":[.....]
			}
		- Monta o prompt usando o template + metadata
		A aplicação deve ser capaz de interpretar tags como <key.key[0]>  ou <key.key[...]> ( para colocar todos os elementos).
		Pensar sobre outras possibilidades que possam ser uteis
		Envia mensagem para a fila do ai-request com o seguintes dados:	*/

	slog.DebugContext(ctx, "useCase.Handle",
		slog.String("details", "process finished"))
	return nil
}

func NewUseCase() interfaces.UseCase {
	return &UseCase{}
}
